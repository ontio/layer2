<template>
	<div class="home">
		<p class="link"
			@click="onBack">
			< {{$t('home.back')}}</p>
				<el-tabs
				v-model="activeName"
				@tab-click="handleClick">
				<el-tab-pane :label="$t('home.transferRecords')"
					name="first">
					<el-table :data="txRecords">
						<el-table-column prop="FromAddress"
							label="From"></el-table-column>
						<el-table-column prop="ToAddress"
							label="To"></el-table-column>
						<el-table-column prop="Token"
							label="Token"></el-table-column>
						<el-table-column prop="Amount"
							label="Amount"></el-table-column>
						<el-table-column prop="TxHash"
							label="TxHash"></el-table-column>
						<el-table-column prop="Date"
							label="Date"></el-table-column>
					</el-table>
				</el-tab-pane>
				<el-tab-pane :label="$t('home.depositRecords')"
					name="second">
                    <el-table :data="depositRecords">
						<el-table-column prop="FromAddress"
							label="From"></el-table-column>
						<el-table-column prop="ToAddress"
							label="To"></el-table-column>
						<el-table-column prop="Token"
							label="Token"></el-table-column>
						<el-table-column prop="Amount"
							label="Amount"></el-table-column>
						<el-table-column prop="TxHash"
							label="TxHash"></el-table-column>
						<el-table-column prop="Date"
							label="Date"></el-table-column>
					</el-table>
                </el-tab-pane>
				<el-tab-pane :label="$t('home.withdrawRecords')"
					name="third">
                    <el-table :data="withdrawRecords">
						<el-table-column prop="FromAddress"
							label="From"></el-table-column>
						<el-table-column prop="ToAddress"
							label="To"></el-table-column>
						<el-table-column prop="Token"
							label="Token"></el-table-column>
						<el-table-column prop="Amount"
							label="Amount"></el-table-column>
						<el-table-column prop="TxHash"
							label="TxHash"></el-table-column>
						<el-table-column prop="Date"
							label="Date"></el-table-column>
					</el-table>
                    </el-tab-pane>
				</el-tabs>
	</div>
</template>
<script>
import dateFormat from "dateformat";
import { mapState } from "vuex";
import BigNumber from "bignumber.js";

export default {
	data() {
		return {
			activeName: "first",
            txRecords: [],
            depositRecords: [],
            withdrawRecords: []
		};
	},
	computed: {
		...mapState({
			onchain_tokens: state => state.onchain_tokens,
			offchain_tokens: state => state.offchain_tokens,
			address: state => state.address,
			privateKey: state => state.privateKey
		})
	},
	mounted() {
        this.fetchTxRecords();
        this.fetchDepositRecords();
        this.fetchWithdrawRecords()
	},
	methods: {
		onBack() {
			this.$router.push("/");
		},
		handleClick(tab, event) {
		},
		async fetchTxRecords() {
			const url =
				process.env.VUE_APP_LAYER2_SERVER + "/api/v1/getlayer2tx/" + this.address;
			const res = await fetch(url, { method: "get" });
			const result = await res.json();
			console.log(result);
			if (result.code === 1) {
				const records = JSON.parse(result.result);
				records.forEach(item => {
					item.Date = dateFormat(
						new Date(Number(item.TT) * 1000),
						"yyyy-mm-dd HH:MM:ss"
					);
					item.Token =
						item.TokenAddress ===
						"0000000000000000000000000000000000000002"
							? "XONG"
							: "XONT";
					if (item.Token === "XONG") {
						item.Amount = new BigNumber(item.Amount)
							.dividedBy(1e9)
							.toNumber();
					}
				});
				this.txRecords = records;
			}
        },
        async fetchWithdrawRecords() {
			const url =
				process.env.VUE_APP_LAYER2_SERVER + "/api/v1/getlayer2withdraw/" + this.address;
			const res = await fetch(url, { method: "get" });
			const result = await res.json();
			console.log(result);
			if (result.code === 1) {
				const records = JSON.parse(result.result);
				records.forEach(item => {
					item.Date = dateFormat(
						new Date(Number(item.TT) * 1000),
						"yyyy-mm-dd HH:MM:ss"
					);
					item.Token =
						item.TokenAddress ===
						"0000000000000000000000000000000000000002"
							? "XONG"
							: "XONT";
					if (item.Token === "XONG") {
						item.Amount = new BigNumber(item.Amount)
							.dividedBy(1e9)
							.toNumber();
					}
				});
				this.withdrawRecords = records;
			}
        },
        async fetchDepositRecords() {
			const url =
				process.env.VUE_APP_LAYER2_SERVER + "/api/v1/getlayer2deposit/" + this.address;
			const res = await fetch(url, { method: "get" });
			const result = await res.json();
			console.log(result);
			if (result.code === 1) {
				const records = JSON.parse(result.result);
				records.forEach(item => {
					item.Date = dateFormat(
						new Date(Number(item.TT) * 1000),
						"yyyy-mm-dd HH:MM:ss"
					);
					item.Token =
						item.TokenAddress ===
						"0000000000000000000000000000000000000002"
							? "XONG"
							: "XONT";
					if (item.Token === "XONG") {
						item.Amount = new BigNumber(item.Amount)
							.dividedBy(1e9)
							.toNumber();
					}
				});
				this.depositRecords = records;
			}
		}
	}
};
</script>
<style lang="scss" scoped>
.link {
	color: rgba(21, 44, 255, 1);
	cursor: pointer;
	text-align: left;
	margin-bottom: 20px;
}
.link:hover {
	opacity: 0.8;
}
</style>