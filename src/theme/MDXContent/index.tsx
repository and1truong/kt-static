import React, {type ReactNode} from 'react';
import MDXContent from '@theme-original/MDXContent';
import type MDXContentType from '@theme/MDXContent';
import type {WrapperProps} from '@docusaurus/types';
import {DiscussionEmbed} from 'disqus-react';
import BrowserOnly from "@docusaurus/BrowserOnly";

type Props = WrapperProps<typeof MDXContentType>;

export default function MDXContentWrapper(props: Props): ReactNode {
    return (
        <>
            <MDXContent {...props} />
            
            <BrowserOnly fallback={<div/>}>
                {
                    () => <DiscussionEmbed
                        shortname='https-thanhkinh-vercel-app'
                        config={
                            {
                                identifier: props.children["_owner"]["memoizedProps"].route.path,
                                title: props.children["type"].contentTitle,
                                language: 'vi',
                            }
                        }
                    />
                }
            </BrowserOnly>
        </>
    );
}
