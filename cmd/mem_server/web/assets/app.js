Vue.createApp({
    data() {
        return {
            player_name: "Riandy",
            state: "NEW_GAME",
            question: "1 + 2",
            choices: ["1", "2", "3"],
            score: 0,
            timeout: 1000, // 5 seconds
            duration: 1000, // 5 seconds
            correct_idx: 2,
            questInterval: null
        }
    },
    methods: {
        play() {
            this.state = "NEW_QUESTION"
            this.timeout = 1000;
            this.duration = 1000;
            this.questInterval = setInterval(() => {
                this.duration-=1;
                if (this.duration === 0) {
                    this.submitAnswer(-1);
                }
            }, 5)
        },
        submitAnswer(idx) {
            clearTimeout(this.questInterval);
            this.state = "GAME_OVER";
        }
    }
}).mount('#app')