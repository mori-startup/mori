<template>
    <button class="btn" @click="toggleModal">New group<i class="uil uil-plus"></i></button>

    <Modal v-show="isOpen" @closeModal="toggleModal(); toggleClearInput();">
        <template #title>
            Create new group
        </template>

        <template #body>
            <form @submit.prevent="submitNewGroup" ref="theForm">
                <div class="form-input">
                    <label for="name">Name</label>
                    <input type="text" name="name" id="name">
                </div>

                <div class="form-input">
                    <label for="description">Description</label>
                    <textarea id="description"
                              name="description"
                              rows="4"
                              cols="50"
                              required
                              placeholder="Describe here"></textarea>
                </div>
    
                <MultiselectDropdown

                                     v-model:checkedOptions="checkedFollowers"
                                     :content="getMyFollowersList"
                                     label-name="Invite users"
                                     placeholder="Select users" />

                <button class="btn form-submit" type="submit">Create</button>
            </form>

        </template>
    </Modal>
</template>


<script>

import Modal from "@/components/Modal.vue";
import MultiselectDropdown from "./MultiselectDropdown.vue";

export default {
  components: {
    Modal,
    MultiselectDropdown,
  },
  data() {
    return {
      checkedFollowers: [],
      isOpen: false,
      clearInput: false,
    };
  },
  created() {
    this.getMyFollowers();
  },
  computed: {
    getMyFollowersList() {
      // Retourne la liste des followers depuis Vuex
      return this.$store.getters.followers;
    },
  },
  methods: {
    async getMyFollowers() {
      // Assurez-vous que cette action remplit correctement le state Vuex
      await this.$store.dispatch("getMyFollowers");
    },
    toggleModal() {
      if (this.isOpen) {
        this.$refs.theForm.reset();
        this.checkedFollowers = [];
      }
      this.isOpen = !this.isOpen;
    },
    toggleClearInput() {
      this.clearInput = !this.clearInput;
    },
    async submitNewGroup(e) {
      const form = e.currentTarget;
      const formData = new FormData(form);
      const formDataObject = Object.fromEntries(formData.entries());
      formDataObject["invitations"] = this.getIds(this.checkedFollowers);

      const response = await fetch("http://localhost:8081/newGroup", {
        method: "POST",
        credentials: "include",
        body: JSON.stringify(formDataObject),
      });

      form.reset();
      this.toggleModal();
      this.toggleClearInput();
      this.$store.dispatch("getUserGroups");
    },
    getIds() {
      const arrOfIDS = [];
      for (const name of this.checkedFollowers) {
        for (const obj of this.$store.state.myFollowers) {
          if (obj.nickname === name.nickname) {
            arrOfIDS.push(obj.id);
          }
        }
      }
      return arrOfIDS;
    },
  },
};
</script>


