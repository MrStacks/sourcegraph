import React, { useEffect, useRef } from 'react'

import { mdiSourceFork, mdiArchive, mdiLock } from '@mdi/js'
import classNames from 'classnames'
import SourceRepositoryIcon from 'mdi-react/SourceRepositoryIcon'

import { highlightNode } from '@sourcegraph/common'
import { displayRepoName } from '@sourcegraph/shared/src/components/RepoLink'
import { getRepoMatchLabel, getRepoMatchUrl, RepositoryMatch } from '@sourcegraph/shared/src/search/stream'
import { useCoreWorkflowImprovementsEnabled } from '@sourcegraph/shared/src/settings/useCoreWorkflowImprovementsEnabled'
import { Icon, Link } from '@sourcegraph/wildcard'

import { LastSyncedIcon } from './LastSyncedIcon'
import { ResultContainer } from './ResultContainer'

import styles from './SearchResult.module.scss'

const REPO_DESCRIPTION_CHAR_LIMIT = 500

export interface RepoSearchResultProps {
    result: RepositoryMatch
    onSelect: () => void
    containerClassName?: string
    as?: React.ElementType
    index: number
}

export const RepoSearchResult: React.FunctionComponent<RepoSearchResultProps> = ({
    result,
    onSelect,
    containerClassName,
    as,
    index,
}) => {
    const [coreWorkflowImprovementsEnabled] = useCoreWorkflowImprovementsEnabled()
    const containerElement = useRef<HTMLDivElement>(null)

    const renderTitle = (): JSX.Element => (
        <div className={styles.title}>
            <span
                className={classNames(
                    'test-search-result-label',
                    styles.titleInner,
                    coreWorkflowImprovementsEnabled && styles.mutedRepoFileLink
                )}
            >
                <Link to={getRepoMatchUrl(result)}>{displayRepoName(getRepoMatchLabel(result))}</Link>
            </span>
        </div>
    )

    const renderBody = (): JSX.Element => (
        <div data-testid="search-repo-result">
            <div className={classNames(styles.searchResultMatch, 'p-2 flex-column')}>
                {result.repoLastFetched && <LastSyncedIcon lastSyncedTime={result.repoLastFetched} />}
                <div className="d-flex align-items-center flex-row">
                    <div className={styles.matchType}>
                        <small>Repository match</small>
                    </div>
                    {result.fork && (
                        <>
                            <div className={styles.divider} />
                            <div>
                                <Icon
                                    aria-label="Forked repository"
                                    className={classNames('flex-shrink-0 text-muted', styles.icon)}
                                    svgPath={mdiSourceFork}
                                />
                            </div>
                            <div>
                                <small>Fork</small>
                            </div>
                        </>
                    )}
                    {result.archived && (
                        <>
                            <div className={styles.divider} />
                            <div>
                                <Icon
                                    aria-label="Archived repository"
                                    className={classNames('flex-shrink-0 text-muted', styles.icon)}
                                    svgPath={mdiArchive}
                                />
                            </div>
                            <div>
                                <small>Archived</small>
                            </div>
                        </>
                    )}
                    {result.private && (
                        <>
                            <div className={styles.divider} />
                            <div>
                                <Icon
                                    aria-label="Private repository"
                                    className={classNames('flex-shrink-0 text-muted', styles.icon)}
                                    svgPath={mdiLock}
                                />
                            </div>
                            <div>
                                <small>Private</small>
                            </div>
                        </>
                    )}
                </div>
                {result.description && (
                    <>
                        <div className={styles.dividerVertical} />
                        <div ref={containerElement}>
                            <small>
                                <em>
                                    {result.description.length > REPO_DESCRIPTION_CHAR_LIMIT
                                        ? result.description.slice(0, REPO_DESCRIPTION_CHAR_LIMIT) + ' ...'
                                        : result.description}
                                </em>
                            </small>
                        </div>
                    </>
                )}
            </div>
        </div>
    )

    useEffect((): void => {
        if (containerElement.current && result.descriptionMatches) {
            const visibleDescription = containerElement.current.querySelector('small em')
            if (visibleDescription) {
                for (const range of result.descriptionMatches) {
                    highlightNode(visibleDescription as HTMLElement, range.start.offset, range.end.offset - range.start.offset)
                }
            }
        }
    }, [result.description, result.descriptionMatches])

    return (
        <ResultContainer
            index={index}
            icon={SourceRepositoryIcon}
            collapsible={false}
            defaultExpanded={true}
            title={renderTitle()}
            resultType={result.type}
            onResultClicked={onSelect}
            expandedChildren={renderBody()}
            repoName={result.repository}
            repoStars={result.repoStars}
            className={containerClassName}
            as={as}
        />
    )
}
