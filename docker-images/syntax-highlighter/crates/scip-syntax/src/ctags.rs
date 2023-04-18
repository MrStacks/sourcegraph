use std::{io::BufWriter, ops::Not, path};

use scip::types::{descriptor::Suffix, Descriptor};
use scip_treesitter_languages::parsers::BundledParser;
use serde::{Deserialize, Serialize};

use crate::{get_globals, globals::Scope};

#[derive(Debug)]
pub enum TagKind {
    Function,
    Class,
}

#[derive(Debug)]
pub struct TagEntry {
    pub descriptors: Vec<Descriptor>,
    pub kind: TagKind,
    pub parent: Option<Box<TagEntry>>,

    pub line: usize,
    // pub column: usize,
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(tag = "command")]
pub enum Request {
    #[serde(rename = "generate-tags")]
    GenerateTags {
        // command == generate-tags
        filename: String,
        size: usize,
    },
}

#[derive(Serialize, Debug)]
#[serde(tag = "_type")]
pub enum Reply<'a> {
    #[serde(rename = "program")]
    Program { name: String, version: String },
    #[serde(rename = "completed")]
    Completed { command: String },
    #[serde(rename = "error")]
    Error { message: String, fatal: bool },
    #[serde(rename = "tag")]
    Tag(Tag<'a>),
}

impl<'a> Reply<'a> {
    pub fn write<W: std::io::Write>(self, writer: &mut W) {
        writer
            .write_all(serde_json::to_string(&self).unwrap().as_bytes())
            .unwrap();
        writer.write_all("\n".as_bytes()).unwrap();
    }

    pub fn write_tag<W: std::io::Write>(
        writer: &mut W,
        scope: &Scope,
        path: &'a str,
        language: &'a str,
        tag_scope: Option<&'a str>,
    ) {
        let descriptors = &scope.descriptors;
        let name = descriptors
            .iter()
            .map(|d| d.name.as_str())
            .collect::<Vec<_>>()
            .join(".");

        let tag = Self::Tag(Tag {
            name,
            path,
            language,
            line: scope.range[0] as usize + 1,
            kind: descriptors_to_kind(&scope.descriptors),
            scope: tag_scope,
        });

        tag.write(writer);
    }
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Tag<'a> {
    name: String,
    path: &'a str,
    language: &'a str,
    /// Starts at 1
    line: usize,
    kind: &'a str,
    scope: Option<&'a str>,
    // Can't find any uses of these. If someone reports a bug, we can support this
    // scope_kind: Option<String>,
    // signature: Option<String>,
}

pub fn generate_tags<W: std::io::Write>(
    buf_writer: &mut BufWriter<W>,
    filename: String,
    file_data: &[u8],
) -> Option<()> {
    let path = path::Path::new(&filename);
    let extension = path.extension()?.to_str()?;
    let filepath = path.file_name()?.to_str()?;

    let parser = BundledParser::get_parser_from_extension(extension)?;
    let (root_scope, _) = match get_globals(parser, file_data)? {
        Ok(vals) => vals,
        Err(err) => {
            // TODO: Not sure I want to keep this or not
            #[cfg(debug_assertions)]
            if true {
                panic!("Could not parse file: {}", err);
            }

            return None;
        }
    };

    emit_tags_for_scope(buf_writer, filepath, vec![], &root_scope, "go");
    Some(())
}

fn descriptors_to_kind(descriptors: &[Descriptor]) -> &'static str {
    match descriptors
        .last()
        .unwrap_or_default()
        .suffix
        .enum_value_or_default()
    {
        Suffix::Namespace => "namespace",
        Suffix::Package => "package",
        Suffix::Method => "method",
        Suffix::Type => "type",
        _ => "variable",
    }
}

fn emit_tags_for_scope<W: std::io::Write>(
    buf_writer: &mut BufWriter<W>,
    path: &str,
    parent_scopes: Vec<String>,
    scope: &Scope,
    language: &str,
) {
    let curr_scopes = {
        let mut curr_scopes = parent_scopes.clone();
        for desc in &scope.descriptors {
            curr_scopes.push(desc.name.clone());
        }
        curr_scopes
    };

    if !scope.descriptors.is_empty() {
        let tag_scope = parent_scopes
            .is_empty()
            .not()
            .then(|| parent_scopes.join("."));
        let tag_scope = tag_scope.as_deref();

        Reply::write_tag(&mut *buf_writer, scope, path, language, tag_scope);
    }

    for subscope in &scope.children {
        emit_tags_for_scope(buf_writer, path, curr_scopes.clone(), subscope, language);
    }

    for global in &scope.globals {
        let mut scope_name = curr_scopes.clone();
        scope_name.extend(
            global
                .descriptors
                .iter()
                .take(global.descriptors.len() - 1)
                .map(|d| d.name.clone()),
        );

        Reply::Tag(Tag {
            name: global.descriptors.last().unwrap().name.clone(),
            path,
            language,
            line: global.range[0] as usize + 1,
            kind: descriptors_to_kind(&global.descriptors),
            scope: scope_name
                .is_empty()
                .not()
                .then(|| scope_name.join("."))
                .as_deref(),
        })
        .write(buf_writer);
    }
}