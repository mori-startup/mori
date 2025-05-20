<template>

    <p class="notificationMsg">
        <!-- who sent -->
        <span class="who"><router-link :to="{ path: `/profile/${notification.user.id}`}">{{ notification.user.nickname }}</router-link></span>

        <!-- additional text from database -->
        <span>{{ notification.content }}</span>        

        <span v-if="notification.type === 'GROUP_INVITE'" class="who">{{ notification.group.name }}</span>
        <span v-if="notification.type === 'GROUP_REQUEST'" class="who"><router-link :to="{ path: `/group/${notification.group.id}`}">{{ notification.group.name }}
</router-link></span>

    </p>

</template>
<script>
export default {
    props: ['notification']
}
</script>


<style scoped>
.who {
    white-space: nowrap;
}
</style>