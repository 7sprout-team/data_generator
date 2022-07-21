<script lang="ts" setup>
import { computed, onMounted, reactive } from 'vue';
import * as APP from '../wailsjs/go/main/App'

const info =reactive({
  party:2,
  inter:10,
  array:{} as any,
  progress:0
})


const getParties = ()=>{
  const result:number[] = []
  for (let i=0;i<info.party;i++){
    result.push(info.array[i*2])
    result.push(info.array[i*2+1])
  }
  return result
}

onMounted(()=>{
  window.setInterval(async ()=>{
    info.progress =await  APP.Progress()
  },100)
})

const submit = ()=>{
  const parties = getParties()
  console.log(parties)
  APP.Generate(info.inter,getParties())
}


</script>

<template>
  <div id="wrapper">

    <div class="item">
      <label class="label" for="party">参与方数量</label>
      <input type="number" id="party" v-model="info.party">
    </div>
    <div class="item">
      <label class="label" for="inter">交集大小</label>
      <input type="number" id="inter" v-model="info.inter">
    </div>
    
    <div v-for="idx in info.party" class="item">
      <label  class="label">{{idx}}行数/特征数</label>
      <input type="number" v-model="info.array[(idx-1)*2]" >
      <input type="number"  v-model="info.array[(idx-1)*2+1]" >
    </div>
    <div class="item">
      <button @click="submit">提交</button> <div>生成进度{{(info.progress*100).toFixed(2)}}%</div>
    </div>
    <p>说明：</p>
    <p>0. 各方的行数都必须大于交集数量</p>
    <p>1. 只有第一个参与方会带有y，为0、1二分类</p>
    <p>2. 特征名不会重复</p>
    <p>3. 所有特征都是随机整数</p>
    <p>4. 特征不是完全随机生成的，先生成1000行随机数据，然后用这1000行数据重复N次</p>
    <p>5. 生成的文件保存位置和打开目录相同，一般就是程序所在目录</p>

  </div>
</template>

<style>
html {
  color: white;
}
#wrapper {
  padding:24px;
  display: flex;
  align-items: flex-start;
  flex-direction: column;
  width: 80%;
  height: 80%;
}
.label {
 width: 100px;
}
input  {
  widows: 50px;
}
.item {
  display: flex;
  margin-bottom: 8px;
}

</style>
