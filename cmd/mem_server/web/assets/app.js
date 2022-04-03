Vue.createApp({
    data() {
        return {
            player_name: "Riandy",
            state: "NEW_GAME",
            question: "1 + 2",
            choices: ["1", "2", "3"],
            score: 0,
            timeout: 5,
            correct_idx: 2
        }
    },
    methods: {
        play() {
            this.state = "NEW_QUESTION"
        },
        submitAnswer(idx) {
            this.state = "GAME_OVER"
        }
    }
}).mount('#app')