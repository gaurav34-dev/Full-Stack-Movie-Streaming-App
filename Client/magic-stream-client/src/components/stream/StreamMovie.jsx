import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import './StreamMovie.css';

/**
 * Extract a clean YouTube video ID from various formats:
 *  - "dQw4w9WgXcQ"                        (raw ID)
 *  - "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
 *  - "https://youtu.be/dQw4w9WgXcQ"
 *  - "https://www.youtube.com/embed/dQw4w9WgXcQ"
 */
function extractYouTubeId(raw) {
  if (!raw) return null;
  const str = raw.trim();

  // Already a bare 11-char ID
  if (/^[A-Za-z0-9_-]{11}$/.test(str)) return str;

  try {
    const url = new URL(str);
    // youtube.com/watch?v=ID
    if (url.searchParams.has('v')) return url.searchParams.get('v');
    // youtu.be/ID or youtube.com/embed/ID
    const parts = url.pathname.split('/').filter(Boolean);
    const last = parts[parts.length - 1];
    if (last && /^[A-Za-z0-9_-]{11}$/.test(last)) return last;
  } catch {
    // not a URL — try regex on the raw string
    const match = str.match(/(?:v=|\/embed\/|youtu\.be\/)([A-Za-z0-9_-]{11})/);
    if (match) return match[1];
  }

  // Fallback: use whatever was passed (might work, might not)
  return str;
}

const StreamMovie = () => {
  const { yt_id } = useParams();
  const videoId = extractYouTubeId(yt_id);
  const [iframeError, setIframeError] = useState(false);

  useEffect(() => {
    setIframeError(false);
  }, [yt_id]);

  if (!yt_id || !videoId) {
    return (
      <div className="stream-fallback">
        <h4>No video selected</h4>
        <p>Please go back and choose a movie to stream.</p>
      </div>
    );
  }

  const embedUrl = `https://www.youtube.com/embed/${videoId}?autoplay=1&rel=0&modestbranding=1`;
  const watchUrl = `https://www.youtube.com/watch?v=${videoId}`;

  return (
    <div className="react-player-container">
      {!iframeError ? (
        <iframe
          src={embedUrl}
          title="Movie Stream"
          className="stream-iframe"
          allow="autoplay; encrypted-media; picture-in-picture; fullscreen"
          allowFullScreen
          onError={() => setIframeError(true)}
        />
      ) : (
        <div className="stream-fallback">
          <h4>Playback unavailable</h4>
          <p>This video can't be embedded. You can watch it directly on YouTube:</p>
          <a
            href={watchUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="btn btn-danger btn-lg"
          >
            ▶ Watch on YouTube
          </a>
        </div>
      )}
    </div>
  );
};

export default StreamMovie;
