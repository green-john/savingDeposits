<template>
    <v-dialog v-model="showModal" max-width="600px">
        <v-card>
            <v-card-title>
                <span class="headline">New User</span>
            </v-card-title>
            <v-card-text>
                <v-container grid-list-md>
                    <v-layout wrap>
                        <v-flex xs12 sm12 md12>
                            <v-text-field
                                    label="Name"
                                    hint="Name of the user"
                                    v-model.trim="newUserData.name"
                                    persistent-hint
                                    required>
                            </v-text-field>
                        </v-flex>
                        <v-flex xs12 sm12 md12>
                            <v-text-field
                                    label="Role"
                                    hint="User's role"
                                    v-model.trim="newUserData.role"
                                    persistent-hint
                                    required>
                            </v-text-field>
                        </v-flex>
                        <v-flex xs12 sm12 md12>
                            <v-text-field
                                    label="Password"
                                    hint="Password"
                                    :type="'password'"
                                    v-model.trim="newUserData.password"
                                    persistent-hint
                                    required>
                            </v-text-field>
                        </v-flex>
                        <v-flex xs12 sm12 md12>
                            <v-text-field
                                    label="Repeat password"
                                    hint="Repeat password"
                                    :type="'password'"
                                    v-model.trim="newUserData.repeatPassword"
                                    persistent-hint
                                    required>
                            </v-text-field>
                        </v-flex>
                    </v-layout>
                </v-container>
                <small>*indicates required field</small>
            </v-card-text>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="blue darken-1" flat @click="showModal = false">Close</v-btn>
                <v-btn color="blue darken-1" flat @click="createUser()">Save</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script>
    import $users from "./users";

    export default {
        name: "NewApartment",
        props: {
            showModal: Boolean,
        },
        data() {
            return {
                newUserMessage: null,

                newUserData: {
                    name: null,
                    role: null,
                    password: null,
                    repeatPassword: null,
                },
            }
        },
        methods: {
            isPasswordValid() {
                return this.password !== "" && this.password === this.repeatPassword;
            },

            createUser() {
                // If is not set, set it to myself
                if (!this.isPasswordValid()) {
                    throw Error("Passwords must match and be non empty");
                }

                $users.createUser(this.newUserData).then(res => {
                    alert(`User ${res.username} created with id ${res.id}`);
                    this.clearApartmentData();
                    this.showModal = false;
                }).catch(err => {
                    this.newUserMessage = `Error: ${err}`;
                });
            },

            clearApartmentData() {
                this.newUserData = {
                    name: null,
                    role: null,
                    password: null,
                    repeatPassword: null
                };
                this.newUserMessage = "";
            },
        }
    }
</script>

<style scoped>

</style>