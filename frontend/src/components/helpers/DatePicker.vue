<template>
    <v-menu
            v-model="showMenu"
            :close-on-content-click="false"
            :nudge-right="40"
            transition="scale-transition"
            offset-y
            min-width="290px"
    >
        <template v-slot:activator="{ on }">
            <v-text-field
                    v-model="date"
                    :label="label"
                    prepend-icon="mdi-calendar"
                    readonly
                    v-on="on"
            ></v-text-field>
        </template>
        <v-date-picker v-model="date"
                       @input="closeMenu">
        </v-date-picker>
    </v-menu>
</template>

<script>
    export default {
        name: "DatePicker",
        props: {
            value: String,
            label: String,
        },

        data() {
            return {
                date: "",
                showMenu: false
            }
        },

        created() {
            this.date = this.value;
        },

        methods: {
            closeMenu() {
                this.showMenu = false;
                this.$emit('input', this.date);
            }
        }
    }
</script>

<style scoped>

</style>