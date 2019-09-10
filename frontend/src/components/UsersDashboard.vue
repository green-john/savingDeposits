<template>
    <v-data-table
            :headers="headers"
            :items="users"
            sort-by="users"
            class="elevation-1"
    >
        <template v-slot:top>
            <v-toolbar flat color="white">
                <v-toolbar-title>Users</v-toolbar-title>
                <v-divider
                        class="mx-4"
                        inset
                        vertical
                ></v-divider>
                <v-toolbar-title>
                    <router-link to="/dashboard">Dashboard</router-link>
                </v-toolbar-title>
                <v-divider
                        class="mx-4"
                        inset
                        vertical
                ></v-divider>
                <v-toolbar-title>
                    <router-link to="/logout">Log out</router-link>
                </v-toolbar-title>
                <v-divider
                        class="mx-4"
                        inset
                        vertical
                ></v-divider>
                <div class="flex-grow-1"></div>
                <v-dialog v-model="dialog" max-width="500px">
                    <template v-slot:activator="{ on }">
                        <v-btn color="primary" dark class="mb-2" v-on="on">New User</v-btn>
                    </template>
                    <v-card>
                        <v-card-title>
                            <span class="headline">{{ formTitle }}</span>
                        </v-card-title>

                        <v-card-text>
                            <v-container>
                                <v-row>
                                    <v-col cols="12" sm="6" md="4">
                                        <v-text-field v-model="editedUser.username" label="Username"></v-text-field>
                                    </v-col>
                                    <v-col cols="12" sm="6" md="4">
                                        <v-text-field v-model="editedUser.role" label="Role"></v-text-field>
                                    </v-col>
                                    <v-col cols="12" sm="6" md="4">
                                        <v-text-field v-model="editedUser.password" label="Password"></v-text-field>
                                    </v-col>
                                </v-row>
                            </v-container>
                        </v-card-text>

                        <v-card-actions>
                            <div class="flex-grow-1"></div>
                            <v-btn color="blue darken-1" text @click="close">Cancel</v-btn>
                            <v-btn color="blue darken-1" text @click="save">Save</v-btn>
                        </v-card-actions>
                    </v-card>
                </v-dialog>
            </v-toolbar>
        </template>
        <template v-slot:item.action="{ item }">
            <v-icon
                    small
                    class="mr-2"
                    @click="editItem(item)"
            >
                edit
            </v-icon>
            <v-icon
                    small
                    @click="deleteItem(item)"
            >
                delete
            </v-icon>
        </template>
    </v-data-table>
</template>

<script>
    import Vue from "vue";
    import $users from "./users";
    import $auth from "./auth";

    function handleError(err) {
        if (err.response) {
            console.log(err.response.status);
            alert(`[ERROR] ${err.response.data}`);
        } else if (err.request) {
            alert(`[ERROR] ${err.request}`);
        } else {
            alert(`[ERROR] ${err.message}`);
        }

        console.log(err.config);
    }

    export default {
        data: () => ({
            dialog: false,
            headers: [
                {
                    text: 'Username',
                    align: 'left',
                    value: 'username',
                },
                {text: 'Role', value: 'role'},
                {text: 'Actions', value: 'action', sortable: false},
            ],
            users: [],
            editedIndex: -1,
            editedUser: {
                username: '',
                password: '',
                role: '',
            },
            defaultUser: {
                username: '',
                role: '',
            },
        }),
        computed: {
            formTitle() {
                return this.editedIndex === -1 ? 'New User' : 'Edit User'
            },

            loggedIn() {
                return $auth.isLoggedIn();
            },
        },
        watch: {
            dialog(val) {
                val || this.close()
            },
        },

        created() {
            this.getAllUsers();
        },

        methods: {
            assignIncomingUsers(incoming) {
                Vue.set(this.users, 'length', incoming.length);
                for (let i = 0; i < incoming.length; i++) {
                    this.users.splice(i, 1, incoming[i]);
                }
            },

            getAllUsers() {
                $users.getAllUsers().then(res => {
                    console.log(res);
                    this.assignIncomingUsers(res);
                });
            },

            editItem(item) {
                this.editedIndex = this.users.indexOf(item);
                this.editedUser = Object.assign({}, item);
                console.log(this.editedUser);
                this.dialog = true;
            },

            deleteItem(user) {
                if (confirm('Are you sure you want to delete this user?')) {
                    $users.deleteUser(user.id).then(() => {
                        this.getAllUsers();
                    }).catch(err => {
                        handleError(err);
                    });
                }
            },

            close() {
                this.dialog = false;
                setTimeout(() => {
                    this.editedUser = Object.assign({}, this.defaultUser);
                    this.editedIndex = -1
                }, 300)
            },

            save() {
                if (this.editedIndex > -1) {
                    const userToUpdate = this.users[this.editedIndex];
                    console.log(userToUpdate);

                    $users.updateUser(userToUpdate.id, this.editedUser.username, this.editedUser.password,
                        this.editedUser.role).then(res => {
                            alert(`User (id=${res.id}) '${res.username}' updated`);
                            this.getAllUsers();
                            this.close();
                        }
                    ).catch(err => {
                        handleError(err);
                    });
                } else {
                    $users.createUser(this.editedUser.username, this.editedUser.password, this.editedUser.role).then(
                        res => {
                            alert(`User (id=${res.id}) '${res.username}' created`);
                            this.getAllUsers();
                            this.close();
                        }
                    ).catch(err => {
                        handleError(err);
                    });
                }
            },
        },
    }
</script>

<!--<script>-->
<!--    -->
<!--    import $auth from "./auth";-->
<!--    import $users from "./users";-->
<!--    import NewUser from './NewUser';-->

<!--    export default {-->
<!--        name: 'UserDashboard',-->
<!--        components: {NewUser},-->
<!--        data() {-->
<!--            return {-->
<!--                showModal: false,-->
<!--                users: [],-->
<!--                userData: {-->
<!--                    username: null,-->
<!--                    role: null-->
<!--                }-->
<!--            }-->
<!--        },-->

<!--        created() {-->
<!--            this.loadUsers();-->
<!--        },-->

<!--        computed: {-->
<!--            loggedIn() {-->
<!--                return $auth.isLoggedIn();-->
<!--            },-->
<!--        },-->

<!--        methods: {-->
<!--            loadUsers() {-->
<!--                $users.getAllUsers().then(res => {-->
<!--                    this.users = res;-->
<!--                    console.log(this.users);-->
<!--                }).catch(err => {-->
<!--                    alert(err);-->
<!--                });-->
<!--            },-->

<!--            getUserInfo() {-->
<!--                $auth.getUserInfo().then(res => {-->
<!--                    this.userData = res.data;-->
<!--                }).catch(err => {-->
<!--                    alert(err);-->
<!--                })-->
<!--            },-->
<!--        },-->
<!--    }-->
<!--</script>-->

<!--<style scoped>-->
<!--    .dashboard {-->
<!--        display: grid;-->
<!--        height: 100vh;-->
<!--        grid-template-columns: 1fr;-->
<!--        grid-template-rows: 1fr;-->
<!--    }-->

<!--    .userRow {-->
<!--        display: block;-->
<!--    }-->
<!--</style>-->
