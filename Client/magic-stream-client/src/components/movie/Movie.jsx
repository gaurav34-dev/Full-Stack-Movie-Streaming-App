import Button from 'react-bootstrap/Button'
import { Link } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCirclePlay } from '@fortawesome/free-solid-svg-icons';
import "./Movie.css";

/**
 * Extract a clean YouTube video ID so the route stays clean:
 * /stream/dQw4w9WgXcQ  (not /stream/https://youtube.com/...)
 */
function extractYouTubeId(raw) {
  if (!raw) return null;
  const str = raw.trim();
  if (/^[A-Za-z0-9_-]{11}$/.test(str)) return str;
  try {
    const url = new URL(str);
    if (url.searchParams.has('v')) return url.searchParams.get('v');
    const parts = url.pathname.split('/').filter(Boolean);
    const last = parts[parts.length - 1];
    if (last && /^[A-Za-z0-9_-]{11}$/.test(last)) return last;
  } catch {
    const match = str.match(/(?:v=|\/embed\/|youtu\.be\/)([A-Za-z0-9_-]{11})/);
    if (match) return match[1];
  }
  return str;
}

const Movie = ({ movie, updateMovieReview }) => {
  const videoId = extractYouTubeId(movie.youtube_id);

  return (
    <div className="col-md-4 mb-4" key={movie._id}>
      <Link
        to={videoId ? `/stream/${videoId}` : '#'}
        style={{ textDecoration: 'none', color: 'inherit' }}
      >
        <div className="card h-100 shadow-sm movie-card">
          <div style={{ position: "relative" }}>
            <img
              src={movie.poster_path}
              alt={movie.title}
              className="card-img-top"
              style={{
                objectFit: "contain",
                height: "250px",
                width: "100%"
              }}
            />
            <span className="play-icon-overlay">
              <FontAwesomeIcon icon={faCirclePlay} />
            </span>
          </div>
          <div className="card-body d-flex flex-column">
            <h5 className="card-title">{movie.title}</h5>
            <p className="card-text mb-2">{movie.imdb_id}</p>
          </div>
          {movie.ranking?.ranking_name && (
            <span className="badge bg-dark m-3 p-2" style={{ fontSize: "1rem" }}>
              {movie.ranking.ranking_name}
            </span>
          )}
          {updateMovieReview && (
            <Button
              variant="outline-info"
              onClick={e => {
                e.preventDefault();
                updateMovieReview(movie.imdb_id);
              }}
              className="m-3"
            >
              Review
            </Button>
          )}
        </div>
      </Link>
    </div>
  );
};

export default Movie;
