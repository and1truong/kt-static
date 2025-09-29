import React from "react";
import { useAudio } from "../../contexts/AudioContext";
import { Callout } from "@radix-ui/themes";
import { InfoCircledIcon } from "@radix-ui/react-icons";

const AudioPlayer = () => {
  const { audioPaths } = useAudio();

  if (!audioPaths) {
    return <></>;
  }

  return (
    <div>
      <Callout.Root>
        <Callout.Text>
          {audioPaths.map((path, index) => (
            <>
              <div key={index} className={"pb-4 w-full"}>
                <audio controls src={`https://kinhthanh.httlvn.org//${path}`}>
                  Your browser does not support the audio element.
                </audio>
              </div>
            </>
          ))}
        </Callout.Text>
      </Callout.Root>
    </div>
  );
};

export default AudioPlayer;
