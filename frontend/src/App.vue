<template>
  <div>
    <div class="profile-card">
      <h2>Now Playing!</h2>


      <div v-if="nowPlaying.isPlaying" class="song-info">
        <h3>🎶 Currently Playing:</h3>
        <p>Title: {{ nowPlaying.title }}</p>
        <p>Artist: {{ nowPlaying.artist }}</p>
        <p>Album: {{ nowPlaying.album }}</p>
        <p>
          Listen on Spotify:
          <a :href="nowPlaying.url" target="_blank" rel="noopener noreferrer">
            {{ nowPlaying.url }}
          </a>
        </p>
      </div>
      <div v-else>
        <p>No music is currently playing.</p>
      </div>



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
  max-width: 400px;
  margin: auto;
  padding: 20px;
  background-color: #f0f0f0;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  font-family: "Arial Black", Gadget, sans-serif;
}

.song-info {
  margin-top: 20px;
}

a {
  color: #3498db;
  text-decoration: none;
}
</style>
