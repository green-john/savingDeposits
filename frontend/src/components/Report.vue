<template>
    <div>
        <v-toolbar flat v-if="loggedIn">
            <v-toolbar-title>Report</v-toolbar-title>
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
        </v-toolbar>

        <v-container class="pa-0 ml-5 mt-3">
            <v-row align="center" no-gutters>
                <v-col cols="3">
                    <DatePicker v-model="startDate"></DatePicker>
                </v-col>
                <v-col cols="3">
                    <DatePicker v-model="endDate"></DatePicker>
                </v-col>

                <v-col cols="3">
                    <v-btn class="mb-2 ml-2" color="primary" @click="getReport()">Get Report</v-btn>
                </v-col>
            </v-row>
        </v-container>

        <v-container>
            <div class="total">
                Total: ${{totalEarned}}
            </div>
        </v-container>

        <v-container class="deposits">
            <v-card outlined v-for="deposit of deposits" :key="deposit.id">
                <v-row no-gutters class="ml-6 mt-4">
                    <v-col sm="2" md="2" class="overline">From {{ deposit.startDate }}</v-col>
                    <v-col sm="2" md="2" class="overline">To {{ deposit.endDate }}</v-col>
                </v-row>

                <v-row no-gutters class="ml-6 mb-4 mt-2">
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
                                    <label class="amount">Initial Amount: ${{ deposit.initialAmount }}</label>
                                    <label class="earned pl-4">{{ deposit.yearlyInterest }}%▲</label>
                                    <label class="paid pl-4">{{ deposit.tax }}%▼</label>
                                </v-col>
                            </v-row>
                        </v-container>
                    </v-col>

                    <v-col>
                        <v-container>
                            <v-row no-gutters>
                                <v-col class="amount">Earnings: ${{deposit.totalRevenue}}</v-col>
                            </v-row>
                            <v-row no-gutters>
                                <v-col class="amount">Tax Paid: ${{deposit.totalTax}}</v-col>
                            </v-row>
                            <v-row no-gutters>
                                <v-col v-bind:class="{
                                earned: positiveBalance(deposit.totalProfit),
                                paid: negativeBalance(deposit.totalProfit)
                                }">Total:
                                    ${{deposit.totalProfit}}
                                </v-col>
                            </v-row>
                        </v-container>
                    </v-col>
                </v-row>
            </v-card>
        </v-container>

    </div>
</template>


<script>
    import DatePicker from "./helpers/DatePicker";
    import $auth from "./auth";
    import $deposits from './deposits';
    import {handleError} from './http';

    export default {
        name: 'Dashboard',
        components: {DatePicker},
        data() {
            return {
                deposits: [],
                userData: {},

                startDate: "",
                endDate: "",
            }
        },

        created() {
             $auth.getUserInfo().then(res => {
                 this.userData = res.data;
                 console.log(this.userData);
             });
        },

        computed: {
            loggedIn() {
                return $auth.isLoggedIn();
            },

            isAdmin() {
                return this.userData.role === "admin";
            },

            totalEarned() {
                let total = 0;
                for (let deposit of this.deposits) {
                    total += deposit.totalProfit;
                }

                return total;
            }
        },

        methods: {
            positiveBalance(balance) {
                return balance > 0;
            },

            negativeBalance(balance) {
                return balance < 0;
            },

            getReport() {
                $deposits.getReport(this.startDate, this.EndDate).then(res => {
                    this.deposits = res;
                }).catch(err => {
                    handleError(err);
                });
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

    .earned {
        color: #5ad136;
        font-size: small;
    }

    .paid {
        color: #ff6e45;
        font-size: small;
    }

    .total {
        font-size: 2rem;
    }
</style>
