import React, {type ReactNode} from 'react';
import MDXContent from '@theme-original/MDXContent';
import type MDXContentType from '@theme/MDXContent';
import type {WrapperProps} from '@docusaurus/types';
import {DiscussionEmbed} from 'disqus-react';
import BrowserOnly from "@docusaurus/BrowserOnly";

type Props = WrapperProps<typeof MDXContentType>;

export default function MDXContentWrapper(props: Props): ReactNode {
    const {id,title } =  props.children["type"].metadata

    return (
        <>
            <MDXContent {...props} />
            
            <BrowserOnly fallback={<div/>}>
                {
                    () => <DiscussionEmbed
                        shortname='https-thanhkinh-vercel-app'
                        config={
                            {
                                identifier: id,
                                title: title,
                                language: 'vi',
                            }
                        }
                    />
                }
            </BrowserOnly>
        </>
    );
}
