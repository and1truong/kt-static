import React, {type ReactNode} from 'react';
import MDXContent from '@theme-original/MDXContent';
import type MDXContentType from '@theme/MDXContent';
import type {WrapperProps} from '@docusaurus/types';
import {DiscussionEmbed} from 'disqus-react';

type Props = WrapperProps<typeof MDXContentType>;

export default function MDXContentWrapper(props: Props): ReactNode {
    const identifier = props.children["_owner"]["memoizedProps"].route.path
    const pageTitle = props.children["type"].contentTitle

    console.log({ identifier, pageTitle })

    return (
        <>
            <MDXContent {...props} />

            {
                identifier && <>
                    <DiscussionEmbed
                        shortname='https-thanhkinh-vercel-app'
                        config={
                            {
                                identifier,
                                title: pageTitle,
                                language: 'vi',
                            }
                        }
                    />
                </>
            }
        </>
    );
}
