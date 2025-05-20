<template>

    <div class="register__container">
        <Register />
    </div>

</template>

<script>
import Register from "@/components/Register.vue";
export default {
    name: 'RegisterView',
    components: { Register, },
}
</script>


<style>
.register__container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100vh;
}
</style>