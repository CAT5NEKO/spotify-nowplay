<template>
  <div class="profile-card">
    <h2>Now Playing</h2>

    <div v-if="nowPlaying.isPlaying" class="song-info">
      <img
          :src="nowPlaying.albumCoverURL"
          alt="Album Cover"
          class="album-cover"
      />
      <div class="details">
        <p class="title">{{ nowPlaying.title }}</p>
        <p class="artist">{{ nowPlaying.artist }}</p>
        <p class="album">{{ nowPlaying.album }}</p>
        <p class="listen">
          <a :href="nowPlaying.url" target="_blank" rel="noopener noreferrer">
            Listen on Spotify
          </a>
        </p>
      </div>
    </div>
    <div v-else>
      <p>No music is currently playing.</p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      nowPlaying: {
        title: '',
        artist: '',
        album: '',
        url: '',
        isPlaying: false,
      },
      socket: null,
    };
  },
  mounted() {
    this.socket = new WebSocket('ws://localhost:4400/ws');

    this.socket.addEventListener('open', (event) => {
      console.log('WebSocket connection opened:', event);
    });

    this.socket.addEventListener('message', (event) => {
      const data = JSON.parse(event.data);
      this.nowPlaying = data;
    });

    this.socket.addEventListener('close', (event) => {
      console.log('WebSocket connection closed:', event);
    });
  },
  beforeUnmount() {
    if (this.socket) {
      this.socket.close();
    }
  },
};
</script>

<style scoped>
.profile-card {
  max-width: 300px;
  margin: auto;
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  font-family: "Arial", sans-serif;
  text-align: center;
}

h2 {
  font-size: 1.5rem;
  margin-bottom: 20px;
}

.song-info {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
}

.album-cover {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.details {
  text-align: left;
}

.title {
  font-size: 1.2rem;
  font-weight: bold;
}

.artist, .album {
  font-size: 1rem;
  color: #777;
}

.listen {
  margin-top: 10px;
}

a {
  color: #3498db;
  text-decoration: none;
}
</style>
