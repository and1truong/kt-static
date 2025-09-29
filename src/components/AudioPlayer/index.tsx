import React from "react";
import { useAudio } from "../../contexts/AudioContext";

const AudioPlayer = () => {
    const { audioPaths } = useAudio();

    if (!audioPaths) {
        return <></>
    }

    return (
        <div>
            {
                audioPaths.map( (path, index ) => <>
                    <div key={index}>
                        <audio controls src={`https://kinhthanh.httlvn.org//${ path }`}>
                            Your browser does not support the audio element.
                        </audio>
                    </div>
                </>)
            }
        </div>
    );
};

export default AudioPlayer;
