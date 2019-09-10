<template>
    <div class="wrapper">
        <v-toolbar flat>
            <v-toolbar-title>
                Dashboard
                <router-link class="logout-link" v-if="loggedIn" to="/logout">Log out</router-link>
            </v-toolbar-title>
            <div class="flex-grow-1"></div>
            <v-dialog v-on:click:outside="close()" v-model="showModal" max-width="500px"
                      v-on:keydown.esc="close()">
                <template v-slot:activator="{ on }">
                    <v-btn color="primary" dark class="mb-2" v-on="on">New Deposit</v-btn>
                </template>
                <v-card class="pa-0">
                    <v-card-title>
                        <span class="headline">{{ fromTitle }}</span>
                    </v-card-title>

                    <v-card-text>
                        <v-container>
                            <v-row>
                                <v-col cols="12" sm="6" md="4">
                                    <v-text-field v-model="editedDeposit.accountNumber"
                                                  label="Acc Number"></v-text-field>
                                </v-col>
                                <v-col cols="12" sm="6" md="4">
                                    <v-text-field v-model="editedDeposit.bankName" label="Bank Name"></v-text-field>
                                </v-col>
                                <v-col cols="12" sm="6" md="4">
                                    <v-text-field v-model.number="editedDeposit.initialAmount"
                                                  label="Initial amount" type="number"></v-text-field>
                                </v-col>
                                <v-col cols="12" sm="6" md="4">
                                    <v-menu
                                            v-model="showMenu.startDate"
                                            :close-on-content-click="false"
                                            :nudge-right="40"
                                            transition="scale-transition"
                                            offset-y
                                            min-width="290px"
                                    >
                                        <template v-slot:activator="{ on }">
                                            <v-text-field
                                                    v-model="editedDeposit.startDate"
                                                    label="Start date"
                                                    prepend-icon="event"
                                                    readonly
                                                    v-on="on"
                                            ></v-text-field>
                                        </template>
                                        <v-date-picker v-model="editedDeposit.startDate"
                                                       @input="showMenu.startDate = false">

                                        </v-date-picker>
                                    </v-menu>
                                </v-col>
                                <v-col cols="12" sm="6" md="4">
                                    <v-menu
                                            v-model="showMenu.endDate"
                                            :close-on-content-click="false"
                                            :nudge-right="40"
                                            transition="scale-transition"
                                            offset-y
                                            min-width="290px"
                                    >
                                        <template v-slot:activator="{ on }">
                                            <v-text-field
                                                    v-model="editedDeposit.endDate"
                                                    label="End Date"
                                                    prepend-icon="event"
                                                    readonly
                                                    v-on="on"
                                            ></v-text-field>
                                        </template>
                                        <v-date-picker v-model="editedDeposit.endDate"
                                                       @input="showMenu.endDate = false">
                                        </v-date-picker>
                                    </v-menu>
                                </v-col>
                                <v-col cols="12" sm="6" md="4">
                                    <v-text-field v-model.number="editedDeposit.yearlyInterest"
                                                  label="Yearly Interest" type="number"></v-text-field>
                                </v-col>
                                <v-col cols="12" sm="6" md="4">
                                    <v-text-field v-model.number="editedDeposit.yearlyTax"
                                                  label="Yearly Tax" type="number"></v-text-field>
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


        <div class="deposits">
            <v-card outlined v-for="deposit of deposits" :key="deposit.id">
                <v-row no-gutters class="ml-6 mt-4">
                    <v-col sm="2" md="2" class="overline">From {{ deposit.startDate }}</v-col>
                    <v-col sm="2" md="2" class="overline">To {{ deposit.endDate }}</v-col>
                </v-row>

                <v-row no-gutters class="ml-6 mt-4">
                    <v-col sm="4" md="3">
                        <v-container>
                            <v-row no-gutters>
                                <v-col class="accountNumber">Account: {{deposit.accountNumber}}</v-col>
                            </v-row>
                            <v-row no-gutters>
                                <v-col class="bankName">Bank: {{deposit.bankName}}</v-col>
                            </v-row>
                            <v-row no-gutters>
                                <v-col>
                                    <label class="initialAmount">Amount: ${{ deposit.initialAmount }}</label>
                                    <label class="interest pl-4">{{ deposit.yearlyInterest }}%▲</label>
                                    <label class="tax pl-4">{{ deposit.yearlyTax }}%▼</label>
                                </v-col>
                            </v-row>
                        </v-container>
                    </v-col>
                </v-row>

                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn outlined text @click="editDeposit(deposit)">Edit</v-btn>
                    <v-btn outlined text @click="deleteDeposit(deposit.id)">Delete</v-btn>
                </v-card-actions>
            </v-card>
        </div>


    </div>
</template>


<script>
    import Vue from 'vue'
    import $auth from "./auth";
    import $deposits from './deposits';
    import $users from "./users";

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
        name: 'Dashboard',
        data() {
            return {
                deposits: [],
                users: [],

                showModal: false,
                showMenu: {
                    endDate: false,
                    startDate: false
                },

                editedIndex: -1,
                editedDeposit: {
                    bankName: '',
                    accountNumber: '',
                    initialAmount: 0.0,
                    yearlyInterest: 0.0,
                    yearlyTax: 0.0,
                    startDate: null,
                    endDate: null,
                    ownerId: 1,
                },

                defaultDeposit: {
                    bankName: '',
                    accountNumber: '',
                    initialAmount: 0.0,
                    yearlyInterest: 0.0,
                    yearlyTax: 0.0,
                    startDate: new Date(2018, 4, 20).toISOString().substr(0, 10),
                    endDate: new Date(2018, 4, 21).toISOString().substr(0, 10),
                    ownerId: 1,
                },
            }
        },

        computed: {
            loggedIn() {
                return $auth.isLoggedIn();
            },

            fromTitle() {
                return this.editedIndex === -1 ? "New User" : "Edit User";
            },
        },

        watch: {
            dialog(val) {
                val || this.close()
            },
        },

        created() {
            this.getAllDeposits();
            this.getUserInfo();
            this.tryLoadingAllUsers();
        },

        methods: {
            getAllDeposits() {
                $deposits.loadAllDeposits().then(res => {
                    console.log(res);
                    this.assignIncomingDeposits(res);
                    console.log(this.deposits);
                }).catch(err => {
                    alert(err);
                });
            },

            assignIncomingDeposits(incoming) {
                Vue.set(this.deposits, 'length', incoming.length);
                for (let i = 0; i < incoming.length; i++) {
                    // this.deposits.$set(i, incoming[i]);
                    this.deposits.splice(i, 1, incoming[i]);
                }
            },

            getUserInfo() {
                $auth.getUserInfo().then(res => {
                    this.userData = res.data;
                }).catch(err => {
                    alert(err);
                })
            },

            tryLoadingAllUsers() {
                $users.getAllUsers().then(res => {
                    this.allUsers = res;
                }).catch(err => {
                    console.log(err);
                })
            },

            close() {
                this.showModal = false;
                setTimeout(() => {
                    this.editedDeposit = Object.assign({}, this.defaultDeposit);
                    this.editedIndex = -1;
                }, 300)
            },

            editDeposit(deposit) {
                this.editedIndex = this.deposits.indexOf(deposit);
                this.editedDeposit = Object.assign({}, deposit);
                console.log(this.editedDeposit);
                this.showModal = true;
            },

            save() {
                if (this.editedIndex > -1) {
                    const depositToUpdate = this.deposits[this.editedIndex];
                    $deposits.updateDeposit(depositToUpdate.id, this.editedDeposit).then(res => {
                            alert(`Deposit updated`);
                            this.getAllDeposits();
                            this.close();
                        }
                    ).catch(err => {
                        handleError(err);
                    });
                } else {
                    console.log(this.editedDeposit);
                    $deposits.createDeposit(this.editedDeposit).then(
                        _ => {
                            alert(`Deposit created`);
                            this.getAllDeposits();
                            this.close();
                        }
                    ).catch(err => {
                        handleError(err);
                    });
                }
            },

            deleteDeposit(depositId) {
                $deposits.deleteDeposit(depositId).then(() => {
                    this.getAllDeposits();
                }).catch(err => {
                    alert(err)
                });
            },
        },
    }
</script>

<style scoped>
    .logout-link {
        position: absolute;
        top: 18px;
        left: 170px;
    }

    .accountNumber {
        font-size: x-large;
    }

    .bankName {
        color: gray;
        font-size: small;
        text-transform: uppercase;
    }

    .initialAmount {
        color: gray;
        font-size: small;
    }

    .interest {
        color: #5ad136;
        font-size: small;
    }

    .tax {
        color: #ff6e45;
        font-size: small;
    }
</style>
