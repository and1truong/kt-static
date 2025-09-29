import React, {useState, type ReactNode} from 'react';
import Heading from '@theme/Heading';
import type {Props} from '@theme/MDXComponents/Heading';
import {useAudio} from "../../contexts/AudioContext";
import AudioPlayer from "../../components/AudioPlayer";
import {MdSpatialAudioOff} from "react-icons/md";

export default function MDXHeading(props: Props): ReactNode {
    const audio = useAudio()
    const [showAudio, setShowAudio] = useState(false)

    if (props.as === "h1") {
        if (audio.audioPaths) {
            return <>
                <h1 {...props}>
                    {props.children}
                    <button
                        className={"audio-toggle"}
                        onClick={() => setShowAudio(!showAudio)}
                        style={{ opacity: showAudio ? 0.4 : 1 }}
                        aria-label={showAudio ? "Hide audio player" : "Show audio player"}
                    >
                        <MdSpatialAudioOff/>
                    </button>
                </h1>
                {showAudio && <AudioPlayer/>}
            </>
        }
    }

    return <Heading {...props} />;
}
