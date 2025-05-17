<template>
  <div class="song-info">
    <img
        :src="track.albumCoverURL"
        alt="Album Cover"
        class="album-cover"
    />
    <div class="details">
      <p class="title">{{ track.title }}</p>
      <p class="artist">{{ track.artist }}</p>
      <p class="album">{{ track.album }}</p>
      <p class="listen">
        <a :href="track.url" target="_blank" rel="noopener noreferrer">
          ワイもSpotifyで聞くで
        </a>
      </p>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    track: {
      type: Object,
      required: true
    }
  }
};
</script>

<style scoped>
.song-info {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 60px;
  flex-wrap: wrap;
}

.album-cover {
  width: 300px;
  height: 300px;
  border-radius: 50%;
  object-fit: cover;
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.1);
  animation: spin 10s linear infinite;
  transition: width 0.3s, height 0.3s;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.details {
  text-align: left;
  flex: 1;
  min-width: 250px;
}

.title {
  font-size: 3.1rem;
  font-weight: bold;
  color: antiquewhite;
}

.artist,
.album {
  font-size: 2rem;
  color: #dcd8d8;
}

.listen {
  margin-top: 30px;
}

a {
  font-size: 1rem;
  color: #45db34;
  text-decoration: none;
}

@media (max-width: 768px) {
  .song-info {
    flex-direction: column;
    align-items: center;
    gap: 30px;
  }

  .album-cover {
    width: 200px;
    height: 200px;
  }

  .title {
    font-size: 2.4rem;
    text-align: center;
  }

  .artist,
  .album {
    font-size: 1.6rem;
    text-align: center;
  }

  .details {
    text-align: center;
  }
}
</style>

