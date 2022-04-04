const apiKey = "c4211664-47dc-4887-a2fe-9e694fbaf55a"

function getNowUnix() {
    return Math.floor(Date.now() / 1000);
}

const app = Vue.createApp({
    data() {
        return {
            player_name: "",
            game_id: "",
            scenario: "NEW_GAME",
            question: "",
            choices: [],
            score: 0,
            max_duration: 0,
            duration: 0,
            correct_idx: -1,
            questInterval: null,
            start_at: 0
        }
    },
    methods: {
        async play() {
            // start new game
            let response = await fetch("/games", {
                method: 'POST',
                headers: { 'X-API-Key': apiKey, 'Content-Type': 'application/json' },
                body: JSON.stringify({player_name: this.player_name})
            });
            let data = (await response.json()).data;
            this.game_id = data.game_id;
            this.fetchNewQuestion();
        },
        async fetchNewQuestion() {
            // fetch new question
            response = await fetch(`/games/${this.game_id}/question`, {
                method: 'PUT',
                headers: { 'X-API-Key': apiKey },
            })
            data = (await response.json()).data;
            // update state
            this.question = data.problem;
            this.choices = data.choices;
            this.start_at = getNowUnix();
            this.scenario = data.scenario;
            // start timer
            this.startTimer(data.timeout, () => {
                // sent no answer if timeout
                this.submitAnswer(-1);
            });
        },
        async submitAnswer(idx) {
            // stop timer
            this.stopTimer();
            // submit answer
            const answer_idx = idx + 1 // because according to api doc answer index start from 1
            let response = await fetch(`/games/${this.game_id}/answer`, {
                method: 'PUT',
                headers: { 'X-API-Key': apiKey, 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    answer_idx: answer_idx,
                    start_at: this.start_at,
                    sent_at: getNowUnix(),
                })
            });
            let data = (await response.json()).data;
            console.log(data);
            // update state
            this.correct_idx = data.correct_idx - 1;
            this.score = data.score;
            if (data.scenario === 'NEW_QUESTION') {
                this.fetchNewQuestion();
                return
            }
            this.scenario = data.scenario;
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
}).mount('#app');