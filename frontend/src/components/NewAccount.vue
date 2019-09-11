<template>
    <div class="wrapper">
        <v-container>
            <h2>New Client Account</h2>
            <v-row justify="center" no-gutters>
                <v-col cols="3">
                    <v-text-field v-model="username"
                                  label="Username"></v-text-field>
                </v-col>
            </v-row>
            <v-row justify="center" no-gutters>
                <v-col cols="3">
                    <v-text-field v-model="password"
                                  label="Password"></v-text-field>
                </v-col>
            </v-row>
            <v-row justify="center" no-gutters>
                <v-col cols="3">
                    <v-text-field v-model="passwordVis"
                                  label="Password (confirm)"></v-text-field>
                </v-col>
            </v-row>
            <v-row justify="center" no-gutters>
                <v-col cols="3">
                    <v-btn @click="newAccount()" color="primary">Create Account</v-btn>
                </v-col>
            </v-row>
            <v-row justify="center" no-gutters>
                <v-col class="mt-2" justify="center" cols="3">
                    <router-link to="/login">Login</router-link>
                    <div class="error">
                        {{ errorMsg }}
                    </div>
                </v-col>
            </v-row>
        </v-container>
    </div>
</template>


<script>
    import $users from './users';
    import {handleError} from "./http";

    export default {
        data() {
            return {
                username: null,
                password: null,
                passwordVis: null,
                errorMsg: null
            }
        },

        methods: {
            newAccount() {
                if (this.password !== this.passwordVis) {
                    this.errorMsg = "Passwords must match";
                    return
                }

                $users.createClientAccount(this.username, this.password).then(res => {
                    alert(`User ${res.username} created`);
                    this.$router.replace('/login');
                }).catch(err => {
                    handleError(err)
                    // this.errorMsg = err.response.data;
                })
            }
        }
    }
</script>

<style scoped>
    /*.wrapper {*/
    /*    display: grid;*/
    /*    justify-content: center;*/
    /*    margin-top: 3rem;*/
    /*    text-align: center;*/
    /*}*/

    /*.newAccount > * {*/
    /*    font-size: 1rem;*/
    /*    height: 2rem;*/
    /*}*/
</style>
