<template>
  <div class="upgrade-ia-page">
    <h1 class="mori">Mori <span class="adder">- Upgrade IA</span></h1>
    <div class="upload-section">
      <div
        class="drop-zone"
        @dragover.prevent
        @dragenter.prevent
        @drop.prevent="handleDrop"
        @click="triggerFileSelect"
      >
        <p>Drop your files here or click to select files.</p>
        <input
          type="file"
          multiple
          ref="fileInput"
          @change="handleFileSelect"
          style="display: none"
        />
      </div>
    </div>

    <!-- Filter Options -->
    <div class="filter-options" v-if="files.length">
      <div class="search-bar">
        <h2>Uploaded Files</h2>
        <input
          type="text"
          v-model="searchQuery"
          placeholder="Search files by name..."
        />
      </div>
      <div class="extra-filters">
        <label>
          <p class="labelfilter">Filter by size</p>
          <select v-model="sizeFilter">
            <option value="all">All</option>
            <option value="small">Small (&lt; 1 MB)</option>
            <option value="medium">Medium (1 - 5 MB)</option>
            <option value="large">Large (&gt; 5 MB)</option>
          </select>
        </label>
        <label>
          <p class="labelfilter">Sort by size</p>
          <select v-model="sizeSort">
            <option value="none">None</option>
            <option value="heaviest">Heaviest First</option>
            <option value="lightest">Lightest First</option>
          </select>
        </label>
        <label>
          <p class="labelfilter">Sort by date</p>
          <select v-model="dateSort">
            <option value="newest">Newest First</option>
            <option value="oldest">Oldest First</option>
          </select>
        </label>
      </div>
    </div>

    <div class="uploaded-files" v-if="filteredFiles.length">
      <ul>
        <li v-for="(file, index) in filteredFiles" :key="index">
          <span class="fileInfo">
            {{ file.name }}
            <span class="fileSize">
              ({{ formatFileSize(file.size) }}) {{ file.uploadDate }}
            </span>
          </span>
          <button @click="deleteFile(file.name)">Delete</button>
        </li>
      </ul>
    </div>
    <div class="uploaded-files-empty" v-else>
      <h2>No files found.</h2>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "UpgradeIA",
  data() {
    return {
      files: [], // Array of objects: { name, size, uploadDate }
      allowedExtensions: [".txt", ".pdf", ".md"],
      searchQuery: "",
      sizeFilter: "all",   // "all", "small", "medium", "large"
      dateSort: "newest",  // "newest" or "oldest"
      sizeSort: "none"     // "none", "heaviest", "lightest"
    };
  },
  computed: {
    filteredFiles() {
      let filtered = this.files;

      // Filter by search query.
      if (this.searchQuery) {
        filtered = filtered.filter(file =>
          file.name.toLowerCase().includes(this.searchQuery.toLowerCase())
        );
      }

      // Filter by file size range.
      if (this.sizeFilter !== "all") {
        filtered = filtered.filter(file => {
          const sizeInMB = file.size / (1024 * 1024);
          if (this.sizeFilter === "small") {
            return sizeInMB < 1;
          } else if (this.sizeFilter === "medium") {
            return sizeInMB >= 1 && sizeInMB <= 5;
          } else if (this.sizeFilter === "large") {
            return sizeInMB > 5;
          }
          return true;
        });
      }

      // Sort: If sizeSort is specified, sort by size.
      // Otherwise, sort by upload date.
      let sorted = filtered.slice(); // create a shallow copy
      if (this.sizeSort !== "none") {
        if (this.sizeSort === "heaviest") {
          sorted.sort((a, b) => b.size - a.size);
        } else if (this.sizeSort === "lightest") {
          sorted.sort((a, b) => a.size - b.size);
        }
      } else {
        sorted.sort((a, b) => {
          const dateA = new Date(a.uploadDate);
          const dateB = new Date(b.uploadDate);
          return this.dateSort === "newest" ? dateB - dateA : dateA - dateB;
        });
      }

      return sorted;
    }
  },
  methods: {
    async handleDrop(event) {
      const droppedFiles = event.dataTransfer.files;
      for (let i = 0; i < droppedFiles.length; i++) {
        await this.uploadFile(droppedFiles[i]);
      }
    },
    async handleFileSelect(event) {
      const selectedFiles = event.target.files;
      for (let i = 0; i < selectedFiles.length; i++) {
        await this.uploadFile(selectedFiles[i]);
      }
    },
    triggerFileSelect() {
      this.$refs.fileInput.click();
    },
    async uploadFile(file) {
      // Check file extension (case-insensitive)
      const fileName = file.name.toLowerCase();
      const extension = fileName.substring(fileName.lastIndexOf("."));
      if (!this.allowedExtensions.includes(extension)) {
        this.$toast.open({
          message: `File type not allowed: ${file.name}. Only .txt, .pdf, and .md files are allowed.`,
          type: "warning",
        });
        return;
      }
  
      const formData = new FormData();
      formData.append("files", file);
      try {
        const response = await axios.post(
          "http://localhost:8081/api/upload",
          formData,
          {
            headers: { "Content-Type": "multipart/form-data" },
            withCredentials: true,
          }
        );
        this.$toast.open({
          message: `${file.name} uploaded successfully`,
          type: "success",
        });
        // Refresh file list after a successful upload.
        this.fetchFiles();
      } catch (error) {
        this.$toast.open({
          message: "Upload error: " + error.message,
          type: "error",
        });
      }
    },
    async fetchFiles() {
      try {
        const response = await axios.get("http://localhost:8081/api/files", {
          withCredentials: true,
        });
        // Ensure files is always an array.
        this.files = response.data ? response.data : [];
      } catch (error) {
        this.$toast.open({
          message: "Fetch files error: " + error.message,
          type: "error",
        });
      }
    },
    async deleteFile(filename) {
      try {
        const response = await axios.delete(
          `http://localhost:8081/api/files/${filename}`,
          { withCredentials: true }
        );
        this.$toast.open({
          message: `${filename} successfully deleted`,
          type: "success",
        });
        // Refresh the file list after deletion.
        this.fetchFiles();
      } catch (error) {
        this.$toast.open({
          message: "Delete error: " + error.message,
          type: "error",
        });
      }
    },
    // Format file size into human-readable form.
    formatFileSize(bytes) {
      if (bytes < 1024) return bytes + " B";
      const kb = bytes / 1024;
      if (kb < 1024) return kb.toFixed(2) + " KB";
      const mb = kb / 1024;
      return mb.toFixed(2) + " MB";
    },
  },
  mounted() {
    this.fetchFiles();
  },
};
</script>

<style scoped>
select {
  background-color: var(--purple-color);
  color: var(--color-white);
  width: 130px;
}

select:hover {
  cursor: pointer;
  background-color: var(--hover-color);
  transition: all 0.3s;
}
.labelfilter {
  margin-bottom: 5px;
}
.fileInfo {
  display: flex;
  width: 94%;
  padding: 10px;
  align-items: center;
  justify-content: space-between;
}
.fileSize {
  color: grey;
}
h1 {
  margin: 25px 0;
}
h2 {
  color: var(--color-white);
  font-weight: bold;
}
p {
  color: var(--color-white);
}
.upgrade-ia-page {
  padding: 20px;
}
.upload-section {
  margin-bottom: 20px;
  position: relative;
}
.drop-zone {
  border: 2px dashed #ccc;
  border-radius: 6px;
  padding: 40px;
  text-align: center;
  cursor: pointer;
  transition: background-color 0.3s ease;
}
.drop-zone:hover {
  background-color: #9146bc39;
}
.search-bar {
  margin-bottom: 10px;
}
.search-bar input {
  padding: 8px;
  border-radius: 10px;
  border: 1px solid #ccc;
  margin-top: 14px;
}
.filter-options {
  margin-bottom: 10px;
  display: flex;
  gap: 20px;
  align-items: center;
}
.filter-options label {
  color: var(--color-white);
}
.extra-filters {
  display: flex;
  gap: 15px;
  margin-top: 5px;
}
.extra-filters label {
  color: var(--color-white);
}
.uploaded-files ul,
.uploaded-files-empty {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  flex-direction: column;
  list-style: none;
  padding: 10px;
  width: 100%;
  overflow-y:scroll;
  height: fit-content;
  border-radius: 10px;
  background-color: var(--bg-neutral);
  max-height: 400px;
}
.uploaded-files li {
  width: 95%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 15px 0;
  color: var(--purple-color);
  border-bottom: 1px solid #a3a1a1;
}
.uploaded-files button {
  background-color: #e74c3c;
  color: #fff;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
}
</style>
