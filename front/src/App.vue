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
            ワイもSpotifyで聞くで
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
  max-width: 900px;
  margin: auto;
  padding: 40px;
  background-color: #383737;
  border-radius: 10px;
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.1);
  font-family: "Arial", sans-serif;
  text-align: center;
}

h2 {
  font-size: 3.2rem;
  text-align: left;
  margin-bottom: 60px;
  color: antiquewhite;
  font-family: "Arial",serif;
}

.song-info {
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: "Arial", sans-serif;
  gap: 60px;
}

.album-cover {
  width: 300px;
  height: 300px;
  border-radius: 50%;
  object-fit: cover;
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.1);
  animation: spin 10s linear infinite;
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
}

.title {
  font-size: 3.1rem;
  font-weight: bold;
  color: antiquewhite;
  font-family: "Arial", sans-serif;
}

.artist, .album {
  font-size: 2rem;
  color: #dcd8d8;
  font-family: "Arial", sans-serif;
}

.listen {
  margin-top: 30px;
}

a {
  font-size: 1rem;
  color: #45db34;
  text-decoration: none;
  font-family: "Arial", sans-serif;
}

@media (max-width: 1200px) {
  .profile-card {
    max-width: 750px;
    padding: 50px;
  }
  h2 {
    font-size: 3.5rem;
    margin-bottom: 50px;
  }
  .album-cover {
    width: 250px;
    height: 250px;
  }
  .title {
    font-size: 3rem;
  }
  .artist, .album {
    font-size: 2.5rem;
  }
  .listen a {
    font-size: 2.5rem;
  }
}

@media (max-width: 900px) {
  .profile-card {
    max-width: 650px;
    padding: 40px;
  }
  h2 {
    font-size: 3rem;
    margin-bottom: 40px;
  }
  .album-cover {
    width: 200px;
    height: 200px;
  }
  .title {
    font-size: 2.5rem;
  }
  .artist, .album {
    font-size: 2rem;
  }
  .listen a {
    font-size: 2rem;
  }
}

@media (max-width: 600px) {
  .profile-card {
    max-width: 500px;
    padding: 30px;
  }
  h2 {
    font-size: 2.5rem;
    margin-bottom: 30px;
  }
  .album-cover {
    width: 150px;
    height: 150px;
  }
  .title {
    font-size: 2rem;
  }
  .artist, .album {
    font-size: 1.8rem;
  }
  .listen a {
    font-size: 1.8rem;
  }
}

@media (max-width: 400px) {
  .profile-card {
    max-width: 90%;
    padding: 20px;
  }
  h2 {
    font-size: 1.2rem;
    margin-bottom: 20px;
  }
  .album-cover {
    width: 120px;
    height: 120px;
  }
  .title {
    font-size: 0.9rem;
  }
  .artist, .album {
    font-size: 0.7rem;
  }
  .listen a {
    font-size: 0.5rem;
  }
}
</style>
