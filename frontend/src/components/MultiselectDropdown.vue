<template>
    <div class="dropdown-wrapper">
      <p class="custom-label">{{ labelName }}</p>
  
      <ul class="checkedOptionsList" v-if="dropdownCheckedOptions.length !== 0">
        <li v-for="checkedOption in dropdownCheckedOptions" :key="checkedOption.id">
          {{ checkedOption.nickname }}
        </li>
      </ul>
  
      <div class="dropdown">
        <button type="button" @click="showDropdown = !showDropdown" class="dropdown-button">
          {{ placeholder }}
          <img class="dropdown-arrow" src="../assets/icons/angle-down.svg" alt="" />
        </button>
  
        <ul class="item-list" v-show="showDropdown">
          <li
            v-if="content !== null && content.length !== 0"
            v-for="option in content"
            :key="option.id"
          >
            <input
              type="checkbox"
              :id="option.id"
              :value="option"
              v-model="dropdownCheckedOptions"
            />
            <label :for="option.id">{{ option.nickname }}</label>
          </li>
  
          <p class="additional-info" v-else>No users to show</p>
        </ul>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    props: ["labelName", "placeholder", "content", "clearInput", "checkedOptions"],
    emits: ["update:checkedOptions", "inputCleared"],
    data() {
      return {
        dropdownCheckedOptions: [],
        showDropdown: false,
        clearedDropdown: false,
      };
    },
  
    watch: {
      dropdownCheckedOptions(newValue) {
        if (this.clearedDropdown) {
          return;
        }
        this.$emit("update:checkedOptions", newValue);
      },
  
      checkedOptions(newValue) {
        if (Object.keys(newValue).length === 0) {
          this.dropdownCheckedOptions = []; // clear dropdown
          this.clearedDropdown = true;
          this.showDropdown = false;
        }
      },
    },
  };
  </script>
  
  <style scoped>
  .dropdown-wrapper {
    display: flex;
    flex-direction: column;
    gap: 5px;
  }
  
  .custom-label {
    color: var(--text-primary); /* Couleur du label */
  }
  
  .dropdown {
    background-color: var(--purple-color);
    box-shadow: var(--container-shadow);
    border-radius: 10px;
    transition: var(--hover-box-shadow-transition);
  }
  
  .dropdown-button {
    padding: 7.5px;
    border: none;
    font-family: "Poppins", sans-serif;
    text-align: left;
    color: var(--color-white);
    background-color: var(--purple-color); /* Fond blanc pour le bouton */
    width: 100%;
    min-height: 35px;
    cursor: pointer;
    border-radius: 5px;
  }
  
  .dropdown-button:hover {
    box-shadow: var(--hover-box-shadow);
  }
  
  .dropdown .item-list {
    padding: 7.5px;
    width: 100%;
    background-color: var(--hover-color); /* Fond des options */
    border-radius: 5px;
  }
  
  .checkedOptionsList {
    display: flex;
    gap: 5px;
    padding: 5px 0;
  }
  
  .checkedOptionsList li {
    background-color: var(--purple-color);
    color: var(--color-white); /* Texte blanc */
    border-radius: 5px;
    padding: 5px;
    font-size: 14px;
  }
  
  /* Styles spécifiques pour les checkboxes */
  .item-list input[type="checkbox"] {
    appearance: none;
    -webkit-appearance: none;
    width: 18px;
    height: 18px;
    border: 2px solid var(--text-primary); /* Bordure visible */
    border-radius: 4px;
    background-color: var(--bg-neutral);
    cursor: pointer;
    outline: none;
    margin-right: 10px;
  }
  
  .item-list input[type="checkbox"]:hover {
    border-color: var(--hover-color);
  }
  
  .item-list input[type="checkbox"]:checked {
    background-color: var(--purple-color);
    border: 2px solid var(--purple-color);
    position: relative;
  }
  
  .item-list input[type="checkbox"]:checked::after {
    content: "";
    position: absolute;
    width: 10px;
    height: 10px;
    background-color: var(--color-white); /* Couleur visible pour le check */
    top: 3px;
    left: 3px;
    border-radius: 2px;
  }
  
  .item-list input[type="checkbox"]::before,
  .item-list input[type="checkbox"]::after {
    content: none;
  }
  
  .item-list label {
    color: var(--text-primary); /* Texte visible */
    font-size: 14px;
    cursor: pointer;
  }
  
  .item-list p {
    font-size: 0.85em;
    color: var(--text-secondary); /* Couleur grise pour les messages */
  }
  </style>
  