<template>
	<div class="home">
        <div class="main-header">
            <el-select v-model="lang" placeholder="请选择" size="small" @change="onChangeLang">
                <el-option value="en" label="EN"></el-option>
                <el-option value="zh" label="中文"></el-option>
            </el-select>
        </div>
		<div class="select-wallet">
            <div class="header">
                <el-button type="primary" @click="()=>{selectWallet = true}">{{$t('home.selectWallet')}}</el-button>
                <p @click="toRecords">{{$t('home.txRecords')}}</p>
            </div>
            <el-dialog :visible.sync="selectWallet" :title="$t('home.selectWallet')">
                <el-form  label-width="200px">
                    <el-form-item :label="$t('home.selectWalletFile')">
                        <el-upload class="upload-demo"
                            action=""
                            :http-request="beforeAvatarUpload" 
                            :auto-upload="true"
                            :show-file-list="false">
                            <el-button size="small"
                                type="primary">{{$t('home.clickSelect')}}</el-button>
                            <p class="file-tip">{{fileName}}</p>
                        </el-upload>
                    </el-form-item>
                    <el-form-item :label="$t('home.inputPassword')">
                        <el-input type="password" v-model="password"></el-input>
                    </el-form-item>
                </el-form>
                <span slot="footer" class="dialog-footer">
                    <el-button @click="selectWallet = false">{{$t('home.cancel')}}</el-button>
                    <el-button type="primary" @click="onDecrypt">{{$t('home.ok')}}</el-button>
                </span>
            </el-dialog>
			
			<p class="address"
				v-if="privateKey"><span>{{$t('home.hasSelect')}}:</span> {{address}}</p>
		</div>
		<div class="token-list ">
            <p class="title">{{$t('home.offchainTip')}}</p>
			<div class="token-item off-chain"
				v-for="(item,index) of offchain_tokens"
				:key="index">
				<div class="token-balance"
					@click="onSelect(index, 'off-chain')">
					<div class="logo">
						<span>{{item.token}}</span>
					</div>
					<div class="balance">
						{{item.balance}}
					</div>
				</div>
				<div class="token-options"
					:class="{'active-options': offchainActiveIndex === index}">
					<div class="option-item"
						@click="onReceive">
						<img src="../assets/receive.png"
							alt="">
						<span>{{$t('home.receive')}}</span>
					</div>
					<div class="option-item"
						@click="onSend(item.token)">
						<img src="../assets/send.png"
							alt="">
						<span>{{$t('home.send')}}</span>
					</div>
					<div class="option-item"
						@click="onConvert(item.token)">
						<img src="../assets/convert.png"
							alt="">
						<span>{{$t('home.convert')}}</span>
					</div>
				</div>

			</div>
		</div>
        <div class="token-list ">
            <p class="title">{{$t('home.onchainTip')}}</p>
			<div class="token-item on-chain"
				v-for="(item,index) of onchain_tokens"
				:key="index">
				<div class="token-balance"
					@click="onSelect(index, 'on-chain')">
					<div class="logo">
						<span>{{item.token}}</span>
					</div>
					<div class="balance">
						{{item.balance}}
					</div>
				</div>
				<div class="token-options"
					:class="{'active-options': onchainActiveIndex === index}">
					<div class="option-item"
						@click="onReceive">
						<img src="../assets/receive.png"
							alt="">
						<span>{{$t('home.receive')}}</span>
					</div>
					<div class="option-item"
						@click="onSend(item.token)">
						<img src="../assets/send.png"
							alt="">
						<span>{{$t('home.send')}}</span>
					</div>
					<div class="option-item"
						@click="onConvert(item.token)">
						<img src="../assets/convert.png"
							alt="">
						<span>{{$t('home.convert')}}</span>
					</div>
				</div>

			</div>
		</div>

        <div class="footer" v-if="isTestnet">
            <p class="network">{{$t('home.currentNetwork')}}：Testnet</p>
            <a href="https://developer.ont.io/" target="_blank">{{$t('home.getTestCoins')}}</a>
        </div>


		<el-drawer :title="$t('home.receiveTokens')"
			:visible.sync="drawer"
			direction="btt"
			size="40%">
			<el-row type="flex"
				justify="center">
				<el-col :xs="24"
					:sm="24"
					:lg="8">
					<div class="receive-container">
						<p class="title">Address: </p>
						<p>{{address}}</p>
						<div class="address-qrcode">
							<vue-qrcode :value="address"></vue-qrcode>
						</div>
					</div>
				</el-col>
			</el-row>

		</el-drawer>
		<el-dialog :title="$t('home.decryptWallet')"
			:visible.sync="passwordVisible">
			<el-form label-width="80px">
				<el-form-item label="Password">
					<el-input type="password"
						v-model="password"></el-input>
				</el-form-item>
			</el-form>
			<span slot="footer"
				class="dialog-footer">
				<el-button @click="passwordVisible = false">{{$t('home.cancel')}}</el-button>
				<el-button type="primary"
					@click="onDecrypt">{{$t('home.ok')}}</el-button>
			</span>
		</el-dialog>
		<el-dialog :title="$t('home.sendTokens')"
			:visible.sync="sendVisible">
			<el-form label-width="80px" ref="sendForm"
				:model="sendForm" :rules="sendRules">
				<el-form-item :label="$t('send.token')" prop="token">
					<el-select v-model="sendForm.token"
						style="float: left;">
						<el-option label="ONT"
							value="ONT"></el-option>
						<el-option label="ONG"
							value="ONG"></el-option>
						<el-option label="XONT"
							value="XONT"></el-option>
						<el-option label="XONG"
							value="XONG"></el-option>

					</el-select>
				</el-form-item>
				<el-form-item :label="$t('send.to')" prop="to">
					<el-input v-model="sendForm.to"></el-input>
				</el-form-item>
				<el-form-item :label="$t('send.amount')" prop="amount">
					<el-input type="number" 
						v-model="sendForm.amount"></el-input>
				</el-form-item>
			</el-form>
			<span slot="footer"
				class="dialog-footer">
				<el-button @click="sendVisible = false">{{$t('home.cancel')}}</el-button>
				<el-button type="primary" v-loading="requesting" :disabled="requesting"
					@click="onSendSubmit">{{$t('home.submit')}}</el-button>
			</span>
		</el-dialog>
		<el-dialog :title="$t('home.convertTokens')" 
			:visible.sync="convertVisible">
			<el-form label-width="80px"  ref="convertForm"
				:model="convertForm" :rules="convertRules">
				<el-form-item :label="$t('convert.from')" prop="from">
					<el-select v-model="convertForm.from"
						style="float: left;"
						@change="onSelectConvertFrom">
						<el-option label="ONT"
							value="ONT"></el-option>
						<el-option label="ONG"
							value="ONG"></el-option>
						<el-option label="XONT"
							value="XONT"></el-option>
						<el-option label="XONG"
							value="XONG"></el-option>

					</el-select>
				</el-form-item>
				<el-form-item :label="$t('convert.to')" prop="to">
					<el-input v-model="convertForm.to"
						disabled></el-input>
				</el-form-item>
				<el-form-item :label="$t('convert.amount')" prop="amount">
					<el-input type="number"
						v-model="convertForm.amount"></el-input>
				</el-form-item>
			</el-form>
			<span slot="footer"
				class="dialog-footer">
				<el-button @click="convertVisible = false">{{$t('home.cancel')}}</el-button>
				<el-button type="primary"  v-loading="requesting" :disabled="requesting"
					@click="onConvertSubmit">{{$t('home.submit')}}</el-button>
			</span>
		</el-dialog>
	</div>
</template>

<script>
import VueQrcode from "vue-qrcode";
import { mapState } from "vuex";
import _ from 'lodash'
// import {Crypto, RpcClient} from 'ontology-ts-sdk'
export default {
	name: "Home",
	components: {
		VueQrcode
	},
	data() {
        const that = this;
		return {
            lang: 'en',
            onchainActiveIndex: -1,
            offchainActiveIndex: -1,
			drawer: false,
			sendVisible: false,
			convertVisible: false,
			sendForm: {
				token: "ONT",
				to: "",
                amount: ""
			},
			convertForm: {
				from: "ONT",
				to: "XONT",
				amount: ""
			},
			passwordVisible: false,
			password: "",
            walletJson: '',
            selectWallet: false,
            fileName: '',
            requesting: false,
            sendRules: {
                token: [
                    {required: true, message: this.$t('send.tokenRequired'), trigger: 'blur'},
                ],
                to: [
                    {required: true, message: this.$t('send.toRequired'), trigger: 'blur'},
                ],
                amount: [
                    {required: true, message: this.$t('send.amountRequired'), trigger: 'blur'},
                ]
            },
            convertRules: {
                from: [
                    {required: true, message: this.$t('convert.fromRequired'), trigger: 'blur'},
                ],
                to: [
                    {required: true, message: this.$t('convert.toRequired'), trigger: 'blur'},
                ],
                amount: [
                    {required: true, message: this.$t('convert.amountRequired'), trigger: 'blur'}
                ]
            }
		};
    },
    mounted() {
        this.getBalance()
        this.intId = setInterval(() => {
            this.getBalance()
        }, 5000)
    },
    beforeDestroy(){
        clearInterval(this.intId)
    },
	computed: {
		...mapState({
            onchain_tokens: state => state.onchain_tokens,
            offchain_tokens: state => state.offchain_tokens,
            address: state => state.address,
            privateKey: state => state.privateKey,
            isTestnet: state => state.isTestnet
		})
	},
	methods: {
        toRecords() {
            this.$router.push('/record')
        },
        onChangeLang() {
            this.$i18n.locale = this.lang;
        },
		beforeAvatarUpload(data) {
			//TODO 解密私钥
            this.fileName = data.file.name;
            this.readWalletFile(data.file).then(json => {
                console.log(json)
                try{
                    let walletJson = JSON.parse(json)
                    if(typeof walletJson === 'string') {
                        walletJson = JSON.parse(walletJson)
                    }
                    this.walletJson = walletJson
                }catch(err) {
                    this.fileName = ''
                    this.$message.error('Invalid wallet file')
                }
            })
		},
		onDecrypt() {
            if(!this.password || !this.walletJson) {
                return;
            }
            const scrypt = this.walletJson.scrypt
            const account = this.walletJson['accounts'][0]
            if(!account) {
                this.$message.error('错误的钱包文件。')
                return;
            }
            const pri = new Ont.Crypto.PrivateKey(account.key);
            const salt = Buffer.from(account.salt, 'base64').toString('hex');
            const address = new Ont.Crypto.Address(account.address);
            this.$store.commit('UPDATE_ADDRESS', address.toBase58())
            try {
                const decrypted = pri.decrypt(this.password, address, salt, {
                cost: scrypt.n, // 除以2时间减半
                blockSize: scrypt.r,
                parallel: scrypt.p,
                size: scrypt.dkLen});
                this.privateKey = decrypted
                this.$store.commit('UPDATE_PRIVATEkEY', decrypted.key)
                this.selectWallet = false
                this.getBalance()
            }catch(err) {
                this.$message.error('解密钱包失败')
                return;
            }
        },
		onSelect(index, type) {
            if(type === 'on-chain') {
                if (this.onchainActiveIndex === index) {
                    this.onchainActiveIndex = -1;
                } else {
                    this.onchainActiveIndex = index;
                    this.offchainActiveIndex = -1;
                }
            } else {
                if (this.offchainActiveIndex === index) {
                    this.offchainActiveIndex = -1;
                } else {
                    this.offchainActiveIndex = index;
                    this.onchainActiveIndex = -1;
                }
            }
			
        },

        getBalance() {
            if(this.address) {
                this.$store.dispatch('getBalance', this.address)
            }
        },
        
		onReceive() {
			this.drawer = true;
		},
		onSend(token) {
            this.sendForm.token = token
			this.sendVisible = true;
		},
		onConvert(token) {
            this.convertForm.from = token
            this.onSelectConvertFrom(token)
			this.convertVisible = true;
		},
		onSelectConvertFrom(token) {
			if (token === "ONT") {
				this.convertForm.to = "XONT";
			} else if (token === "XONT") {
				this.convertForm.to = "ONT";
			} else if (token === "ONG") {
				this.convertForm.to = "XONG";
			} else if (token === "XONG") {
				this.convertForm.to = "ONG";
			}
		},
		onSendSubmit: 
            _.debounce(async function(){
                const valid = await this.$refs["sendForm"].validate();
                if(!valid) return;
                if(this.sendForm.token === 'XONG' || this.sendForm.token === 'ONG') {
                    if(Number(this.sendForm.amount) < 0.1) {
                        this.$message({
                            type: 'warning',
                            message: this.$t('send.minOngAmount')
                        })
                        this.requesting = false;
                        return;
                    }
                }
                 this.requesting = true
                const res = await this.$store.dispatch('sendToken', {...this.sendForm})
                    this.sendVisible = false
                    this.requesting = false
                    if(res.Error === 0) {
                        this.$message.success(this.$t('home.transactionSuccess'))    
                    } else {
                        this.$message.error(res.message || this.$t('home.transactionFail'))    
                    }
            }, 200)
        ,
        onConvertSubmit:
            _.debounce(async function() {
                const valid = await this.$refs["convertForm"].validate();
                if(!valid) return;
                if(this.convertForm.from === 'XONG' || this.convertForm.from === 'ONG') {
                    if(Number(this.convertForm.amount) < 0.1) {
                        this.$message({
                            type: 'warning',
                            message: this.$t('send.minOngAmount')
                        })
                        this.requesting = false;
                        return;
                    }
                } 
                let {from, to ,amount} = this.convertForm
                if(from === 'ONT' || from === 'ONG') {
                    this.requesting = true
                    const res = await this.$store.dispatch('deposit', {amount, asset: from})
                        this.convertVisible = false;
                        this.requesting = false
                        if(res.Error === 0) {
                            this.$message.success(this.$t('home.transactionSuccess'))    
                        } else {
                            this.$message.error(res.message || this.$t('home.transactionFail'))    
                        }
                } else if(from === 'XONT' || from === 'XONG') {
                    this.requesting = true
                    const res = await this.$store.dispatch('withdraw', {amount, asset:from})
                        this.convertVisible = false;
                        this.requesting = false
                        if(res.Error === 0) {
                            this.$message.success(this.$t('home.transactionSuccess'))    
                        } else {
                            this.$message.error(res.message || this.$t('home.transactionFail'))    
                        }
                }
            }, 200),
        
	    readWalletFile($file, readerType) {
			return new Promise(function(resolve, reject) {
				var reader = new FileReader();
				reader.onload = function(params) {
					resolve(reader.result);
				};
				reader.onerror = reject;
				if (readerType === "ArrayBuffer") {
					reader.readAsArrayBuffer($file);
				} else if (readerType === "Binary") {
					reader.readAsBinaryString($file);
				} else if (readerType === "DataUrl") {
					reader.readAsDataURL($file);
				} else {
					reader.readAsText($file);
				}
			});
		}
	}
};
</script>
<style lang="scss" scoped>
.home {
	font-size: 16px;
	font-family: PingFangSC-Medium, PingFang SC;
}
.main-header {
    text-align: right;
    div {
        width: 80px;
    }
}
.token-list {
	margin-top: 20px;
    padding: 0 20px;
    .title {
        text-align: left;
        margin-bottom: 10px;
    }
	.token-item {
		color: #ffffff;
		margin-bottom: 15px;
		.token-balance {
			display: flex;
			justify-content: space-between;
			cursor: pointer;
			padding: 20px 30px;
		}
		.logo {
			font-size: 16px;
			font-family: PingFangSC-Medium, PingFang SC;
			font-weight: 500;
		}
		.balance {
			font-size: 16px;
			font-family: DINPro-Medium, DINPro;
			font-weight: 500;
		}
		.token-options {
			display: flex;
			justify-content: space-between;
			padding: 0 30px;
			max-height: 0;
			padding-bottom: 0;
			overflow: hidden;
			transition: all 0.5s;
		}
		.active-options {
			max-height: 100px;
			padding-bottom: 20px;
		}
		.option-item {
			cursor: pointer;
			display: flex;
			flex-direction: column;
			align-items: center;
			border: 1px solid #ffffff;
			border-radius: 5px;
			width: 80px;
			padding-top: 10px;
			padding-bottom: 5px;
			font-size: 12px;
			img {
				width: 22px;
				height: 20px;
				margin-bottom: 3px;
			}
        }
        .option-item:hover {
            opacity: .8;
        }
    }
    .on-chain {
		background: #409eff;
    }
    .off-chain {
		background: rgba(21, 44, 255, 1);

		
    }
}
.upload-demo {
    text-align:left;
}

.receive-container {
	.title {
		text-align: left;
		opacity: 0.6;
	}
}
.select-wallet {
	margin-top: 50px;
	.address {
		margin: 15px;
		text-align: left;
	}
}
.footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px;
}
.file-tip {
        float: right;
    margin-left: 20px;
}
.header {
    padding: 0 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    p {
        text-decoration: underline;
        color: rgba(21, 44, 255, 1);
        cursor:pointer
    }
    p:hover {
        opacity:.8;
    }
}
</style>
