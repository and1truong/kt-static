import React, { useState, type ReactNode } from "react";
import Heading from "@theme/Heading";
import type { Props } from "@theme/MDXComponents/Heading";
import { useAudio } from "../../contexts/AudioContext";
import AudioPlayer from "../../components/AudioPlayer";
import { Flex, Text, Button } from "@radix-ui/themes";

export default function MDXHeading(props: Props): ReactNode {
  const audio = useAudio();
  const [showAudio, setShowAudio] = useState(false);

  if (props.as === "h1") {
    if (audio.audioPaths) {
      return (
        <>
          <h1 {...props}>
            {props.children}
            <Button
              size="3"
              variant="soft"
              className={"audio-toggle"}
              onClick={() => setShowAudio(!showAudio)}
              style={{ opacity: showAudio ? 0.7 : 1 }}
              aria-label={showAudio ? "Ẩn phần âm thanh" : "Hiện phần âm thanh"}
            >
              {showAudio ? "️🎙 Ẩn phần âm thanh" : "🎙 Hiện phần âm thanh"}
            </Button>
          </h1>
          {showAudio && <AudioPlayer />}
        </>
      );
    }
  }

  return <Heading {...props} />;
}
