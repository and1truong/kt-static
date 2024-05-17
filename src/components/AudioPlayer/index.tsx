import React from "react";

const AudioPlayer = ({ path, x }) => {
    return (
        <div>
            <pre>
                X: {JSON.stringify(x)}
            </pre>
            
            <audio controls src={path}>
                Your browser does not support the audio element.
            </audio>
        </div>
    );
};

export default AudioPlayer;
