<template>
    <div class="wrapper">
        <v-toolbar flat v-if="loggedIn">
            <v-toolbar-title>Dashboard</v-toolbar-title>
            <v-divider
                    class="mx-4"
                    inset
                    vertical
            ></v-divider>
            <v-toolbar-title>
                <router-link to="/report">Report</router-link>
            </v-toolbar-title>
            <v-divider
                    class="mx-4"
                    inset
                    vertical
            ></v-divider>
            <v-toolbar-title v-if="isAdmin">
                <router-link to="/users">Users</router-link>
            </v-toolbar-title>
            <v-divider v-if="isAdmin"
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
                                    <DatePicker label="Start Date" v-model="editedDeposit.startDate"></DatePicker>
                                </v-col>
                                <v-col cols="12" sm="6" md="4">
                                    <DatePicker label="End Date" v-model="editedDeposit.endDate"></DatePicker>
                                </v-col>
                                <v-col cols="12" sm="6" md="4">
                                    <v-text-field v-model.number="editedDeposit.yearlyInterest"
                                                  label="Yearly Interest" type="number"
                                                  hint="Use percentage (e.g. 5 for 5%)"></v-text-field>
                                </v-col>
                                <v-col cols="12" sm="6" md="4">
                                    <v-text-field v-model.number="editedDeposit.tax"
                                                  label="Tax" type="number"
                                                  hint="Use percentage (e.g. 5 for 5%)"></v-text-field>
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

        <v-container class="pa-0 ml-5 mt-3">
            <v-row align="center" no-gutters>
                <v-col cols="2">
                    <v-text-field class="ml-2" v-model.number="filters.minAmount"
                                  label="Min Amount" type="number"></v-text-field>
                </v-col>
                <v-col cols="2">
                    <v-text-field class="ml-2" v-model.number="filters.maxAmount"
                                  label="Max Amount" type="number"></v-text-field>
                </v-col>
                <v-col cols="2">
                    <v-text-field class="ml-2" v-model="filters.bankName"
                                  label="Bank Name"></v-text-field>
                </v-col>
                <v-col cols="2">
                    <DatePicker v-model="filters.startDate"></DatePicker>
                </v-col>
                <v-col cols="2">
                    <DatePicker v-model="filters.endDate"></DatePicker>
                </v-col>

                <v-col cols="1">
                    <v-btn class="mb-2 ml-2" @click="filterData()">
                        <v-icon>mdi-filter</v-icon>
                    </v-btn>
                </v-col>

                <v-col cols="1">
                    <v-btn class="mb-2 ml-2" @click="resetData()">Reset</v-btn>
                </v-col>
            </v-row>
        </v-container>

        <v-container class="deposits">
            <v-card outlined v-for="deposit of deposits" :key="deposit.id">
                <v-row no-gutters class="ml-6 mt-4">
                    <v-col sm="2" md="2" class="overline">From {{ deposit.startDate }}</v-col>
                    <v-col sm="2" md="2" class="overline">To {{ deposit.endDate }}</v-col>
                </v-row>

                <v-row no-gutters class="ml-6 mt-4">
                    <v-col>
                        <v-container>
                            <v-row no-gutters>
                                <v-col class="accountNumber">Account: {{deposit.accountNumber}}</v-col>
                            </v-row>
                            <v-row no-gutters>
                                <v-col class="bankName">Bank: {{deposit.bankName}}</v-col>
                            </v-row>
                            <v-row no-gutters>
                                <v-col>
                                    <label class="amount">Amount: ${{ deposit.initialAmount }}</label>
                                    <label class="interest pl-4">{{ deposit.yearlyInterest }}%▲</label>
                                    <label class="tax pl-4">{{ deposit.tax }}%▼</label>
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
        </v-container>

    </div>
</template>


<script>
    import DatePicker from "./helpers/DatePicker";
    import $auth from "./auth";
    import $deposits from './deposits';
    import $users from "./users";
    import {handleError} from "./http";

    export default {
        name: 'Dashboard',
        components: {DatePicker},
        data() {
            return {
                deposits: [],
                users: [],
                userData: {},

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
                    tax: 0.0,
                    startDate: null,
                    endDate: null,
                    ownerId: 1,
                },

                filters: {
                    bankName: '',
                    minAmount: 0.0,
                    maxAmount: 0.0,
                    startDate: "",
                    endDate: ""
                },

                defaultDeposit: {
                    bankName: '',
                    accountNumber: '',
                    initialAmount: 0.0,
                    yearlyInterest: 0.0,
                    tax: 0.0,
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

            isAdmin() {
                return this.userData.role === "admin";
            },

            fromTitle() {
                return this.editedIndex === -1 ? "New Deposit" : "Edit Deposit";
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
            getAllDeposits(filters) {
                filters = filters || {};
                $deposits.loadAllDeposits(filters).then(res => {
                    this.deposits = res;
                }).catch(err => {
                    handleError(err);
                });
            },

            getUserInfo() {
                $auth.getUserInfo().then(res => {
                    this.userData = res.data;
                }).catch(err => {
                    handleError(err);
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
                    $deposits.updateDeposit(depositToUpdate.id, this.editedDeposit).then(() => {
                            alert(`Deposit updated`);
                            this.getAllDeposits();
                            this.close();
                        }
                    ).catch(err => {
                        handleError(err);
                    });
                } else {
                    if (this.userData.role !== "admin") {
                        this.editedDeposit.ownerId = this.userData.id;
                    }
                    $deposits.createDeposit(this.editedDeposit).then(() => {
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
                    handleError(err)
                });
            },

            filterData() {
                this.getAllDeposits(this.filters);
            },

            resetData() {
                this.getAllDeposits();
            },
        },
    }
</script>

<style scoped>
    .accountNumber {
        font-size: x-large;
    }

    .bankName {
        color: gray;
        font-size: small;
        text-transform: uppercase;
    }

    .amount {
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
