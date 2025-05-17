<template>
  <div class="history-container">
    <h3>再生履歴</h3>
    <ul class="history-list">
      <li v-for="(track, index) in history" :key="index" class="history-item">
        <img :src="track.albumCoverURL" alt="Cover" class="thumbnail" />
        <div class="track-info">
          <p class="title">{{ track.title }}</p>
          <p class="artist">{{ track.artist }}</p>
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  name: 'HistoryList',
  props: {
    history: {
      type: Array,
      required: true
    }
  }
};
</script>

<style scoped>
.history-container {
  margin-top: 60px;
  text-align: left;
  color: #f0f0f0;
}

h3 {
  font-size: 2rem;
  margin-bottom: 20px;
  color: antiquewhite;
}

.history-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.history-item {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  gap: 20px;
  padding: 10px;
  border-radius: 8px;
  background-color: #2e2e2e;
}

.thumbnail {
  width: 50px;
  height: 50px;
  object-fit: cover;
  border-radius: 8px;
  flex-shrink: 0;
}

.track-info {
  flex: 1;
  overflow: hidden;
}

.title {
  font-size: 1.2rem;
  font-weight: bold;
  margin: 0;
  word-wrap: break-word;
  overflow-wrap: break-word;
  white-space: normal;
}

.artist {
  font-size: 1rem;
  color: #bbb;
  margin: 0;
  word-wrap: break-word;
  overflow-wrap: break-word;
  white-space: normal;
}

@media (max-width: 600px) {
  .history-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }

  .thumbnail {
    width: 40px;
    height: 40px;
  }

  .title {
    font-size: 1rem;
  }

  .artist {
    font-size: 0.9rem;
  }
}
</style>
