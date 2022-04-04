Vue.createApp({
    data() {
        return {
            player_name: "Riandy",
            state: "NEW_GAME",
            question: "1 + 2",
            choices: ["1", "2", "3"],
            score: 0,
            max_duration: 0,
            duration: 0,
            correct_idx: 2,
            questInterval: null
        }
    },
    methods: {
        play() {
            this.state = "NEW_QUESTION"
            this.startTimer(5, () => {
                this.submitAnswer(0);
            });
        },
        async submitAnswer(idx) {
            this.stopTimer();
            this.state = "GAME_OVER";
        },
        startTimer(seconds, timeoutCallback) {
            // set current duration & max duration
            this.max_duration = seconds * 200; // the multiplier is 200 to make the animation smooth
            this.duration = this.max_duration;
            // start timer animation
            this.questInterval = setInterval(() => {
                this.duration-=2;
                if (this.duration === 0) {
                    // execute callback
                    timeoutCallback();
                }
            }, 5)
        },
        stopTimer() {
            clearTimeout(this.questInterval);
        }

    }
}).mount('#app')