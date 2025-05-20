<template>
    <div class="sign-in__container">
        <SignIn />
    </div>
</template>

<script>

// @ is an alias to /src
import SignIn from '@/components/SignIn.vue';

export default {
    name: 'SignInView',
    components: {
        SignIn,
    }
}

</script>

<style>
.sign-in__container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100vh;

}

.sign-in__container .sign-in__wrapper {
    place-self: center;
}
</style>
