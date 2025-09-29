import React, { createContext, useContext, ReactNode } from 'react';

interface AudioContextType {
  audioPaths: string[] | null;
}

const AudioContext = createContext<AudioContextType | undefined>(undefined);

export function AudioProvider({ children, audioPaths }: { children: ReactNode; audioPaths: string[] | null }) {
  return (
    <AudioContext.Provider value={{ audioPaths }}>
      {children}
    </AudioContext.Provider>
  );
}

export function useAudio() {
  const context = useContext(AudioContext);
  if (context === undefined) {
    throw new Error('useAudio must be used within an AudioProvider');
  }
  return context;
}