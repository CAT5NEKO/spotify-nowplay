<template>
  <div class="profile-card">
    <h2>Now Playing</h2>
    <TrackCard v-if="nowPlaying.isPlaying" :track="nowPlaying" />
    <div v-else>
      <p>No music is currently playing.</p>
    </div>
    <HistoryList :history="history" />
  </div>
</template>

<script>
import TrackCard from './components/TrackCard.vue';
import HistoryList from './components/HistoryList.vue';

export default {
  components: { TrackCard, HistoryList },
  data() {
    return {
      nowPlaying: {
        title: '',
        artist: '',
        album: '',
        url: '',
        albumCoverURL: '',
        isPlaying: false,
        timestamp: 0,
      },
      lastPlayed: null,
      history: [],
      socket: null,
      reconnectAttempts: 0,
      maxReconnectAttempts: 5,
    };
  },

  mounted() {
    this.createWebSocket();
    document.addEventListener('visibilitychange', this.handleVisibilityChange);
    this.history = [];
  },

  beforeUnmount() {
    if (this.socket) {
      this.socket.close();
    }
    document.removeEventListener('visibilitychange', this.handleVisibilityChange);
  },

  methods: {
    createWebSocket() {
      this.socket = new WebSocket('ws://localhost:4400/ws');

      this.socket.addEventListener('open', () => {
        console.log('WebSocket connection opened');
        this.reconnectAttempts = 0;
      });

      this.socket.addEventListener('message', (event) => {
        try {
          const data = JSON.parse(event.data);
          this.nowPlaying = data;

          if (data.isPlaying) {
            this.addToHistory(data);
          }
        } catch (err) {
          console.error('Error parsing WebSocket message:', err);
        }
      });

      this.socket.addEventListener('close', () => {
        this.handleSocketClose();
      });

      this.socket.addEventListener('error', () => {
        this.handleSocketClose();
      });
    },

    addToHistory(track) {
      const isDuplicate = this.history.some(item => item.title === track.title && item.artist === track.artist);

      if (!isDuplicate) {
        this.history.unshift({ ...track, timestamp: Date.now() });
        if (this.history.length > 5) {
          this.history.pop();
        }
      }
    },

    handleSocketClose() {
      if (this.reconnectAttempts < this.maxReconnectAttempts) {
        this.reconnectAttempts++;
        console.log(`Reconnecting... (${this.reconnectAttempts})`);
        setTimeout(() => {
          this.createWebSocket();
        }, 1000);
      } else {
        console.log('Max reconnect attempts reached.');
      }
    },

    handleVisibilityChange() {
      if (document.hidden && this.socket) {
        this.socket.close();
      } else if (!this.socket || this.socket.readyState === WebSocket.CLOSED) {
        this.createWebSocket();
      }
    },
  }
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
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

a {
  font-size: 1rem;
  color: #45db34;
  text-decoration: none;
}

@media (max-width: 768px) {
  h2 {
    font-size: 2.5rem;
    text-align: center;
  }
}
</style>
