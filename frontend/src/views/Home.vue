<template>
    <Header />
    <!-- 小屏幕 -->
    <div class="block md:hidden mx-4">
        <!-- 快捷入口 -->
        <div class="grid grid-cols-2 gap-3 mt-2">
            <router-link to="/class" class="bg-white rounded-lg shadow-sm p-3 flex items-center gap-3 cursor-pointer hover:shadow-md transition-shadow">
                <div class="w-12 h-12 rounded-full bg-gradient-to-r from-emerald-500 to-teal-500 flex items-center justify-center text-white">
                    <i class="fas fa-book-open text-lg"></i>
                </div>
                <div>
                    <h4 class="text-sm font-medium">所有作品</h4>
                    <p class="text-xs text-gray-500">浏览全部小说</p>
                </div>
            </router-link>
            <router-link to="/rank" class="bg-white rounded-lg shadow-sm p-3 flex items-center gap-3 cursor-pointer hover:shadow-md transition-shadow">
                <div class="w-12 h-12 rounded-full bg-gradient-to-r from-orange-500 to-amber-500 flex items-center justify-center text-white">
                    <i class="fas fa-chart-line text-lg"></i>
                </div>
                <div>
                    <h4 class="text-sm font-medium">排行榜</h4>
                    <p class="text-xs text-gray-500">热门书籍排行</p>
                </div>
            </router-link>
        </div>
        
        <!-- 精品推荐 - 轮播效果 -->
        <section class="mb-6 mt-4">
        <h3 class="text-lg font-bold mb-3 border-b-2 border-primary pb-1 inline-block">精品推荐</h3>
        <div class="relative">
            <!-- 轮播容器 -->
            <div class="overflow-hidden rounded-lg">
                <div class="flex transition-transform duration-500 ease-in-out" :style="{ transform: `translateX(-${featuredCurrentIndex * 100}%)` }">
                    <div v-for="(book, index) in featuredBooks" :key="`featured-${index}`" class="w-full flex-shrink-0">
                        <div class="bg-white rounded-lg overflow-hidden shadow-sm cursor-pointer" @click="navigateToBook(book.id)">
                            <img :src="book.cover" :alt="book.title" class="w-full h-40 object-cover" />
                            <div class="p-3">
                                <h4 class="text-sm font-medium line-clamp-2">{{ book.title }}</h4>
                                <p class="text-xs text-gray-500 mt-1 line-clamp-3">{{ book.description }}</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            
            <!-- 轮播指示器 -->
            <div class="flex justify-center mt-2 gap-1">
                <span v-for="(book, index) in featuredBooks" :key="`featured-dot-${index}`"
                    @click="featuredCurrentIndex = index"
                    class="w-2 h-2 rounded-full transition-all cursor-pointer"
                    :class="featuredCurrentIndex === index ? 'bg-primary' : 'bg-gray-300'">
                </span>
            </div>
        </div>
        </section>

        <!-- 热门推荐 -->
        <section class="mb-6">
        <h3 class="text-lg font-bold mb-3 border-b-2 border-primary pb-1 inline-block">热门推荐</h3>
        <div class="space-y-3">
            <div v-for="(book, index) in recommendBooks" :key="`recommend-${index}`" class="flex gap-2 p-2 bg-white rounded-lg shadow-sm cursor-pointer" @click="navigateToBook(book.id)">
            <img :src="book.cover" :alt="book.title" class="w-16 h-24 object-cover rounded" />
            <div class="flex-1">
                <h4 class="text-sm font-medium line-clamp-2">{{ book.title }}</h4>
                <p class="text-xs text-gray-500 mt-1 line-clamp-2">{{ book.description }}</p>
                <p class="text-xs text-gray-400 mt-1">作者：{{ book.author }}</p>
            </div>
            </div>
        </div>
        </section>

        <!-- 最新更新 -->
        <section class="mb-6">
        <div class="flex justify-between items-center mb-3">
            <h3 class="text-lg font-bold">最新更新</h3>
            <a href="#" class="text-primary text-sm hover:underline">更多 ></a>
        </div>
        <div class="bg-white rounded-lg shadow-sm overflow-hidden">
            <ul class="divide-y divide-gray-200">
            <li v-for="(item, index) in latestBooks" :key="index" class="px-3 py-2 flex justify-between items-start cursor-pointer" @click="navigateToBook(item.id)">
                <div class="flex-1">
                <h4 class="text-sm font-medium line-clamp-1">{{ item.title }}</h4>
                <p class="text-xs text-gray-500 mt-1 line-clamp-2">{{ item.latestChapter }}</p>
                <p class="text-xs text-gray-400 mt-1">作者：{{ item.author }}</p>
                </div>
                <span class="text-xs text-orange-500 font-medium">{{ item.updateTime }}</span>
            </li>
            </ul>
        </div>
        </section>
    </div>


    <!-- 大屏幕 -->
    <div class="hidden md:block container mx-auto px-4">
        <!-- 重点推荐区 -->
        <section class="my-10 w-5xl mx-auto">
            <div class="flex flex-col md:flex-row gap-6 md:items-start">
                <!-- 主推荐 -->
                <div class="w-full md:w-2/3 bg-white rounded-xl shadow-sm overflow-hidden card-hover" @click="navigateToBook(currentBook.id)">
                    <div class="flex flex-col md:flex-row">
                        <!-- 左侧：主图 + 小图轮播 -->
                        <div class="w-full md:w-1/3 relative">
                            <!-- 主图 -->
                            <img :src="currentBook.cover" :alt="currentBook.title" class="w-full h-full object-cover cursor-pointer" />
                            <div v-if="currentBook.source"
                                class="absolute top-2 left-2 bg-red-500 text-white text-xs px-2 py-1 rounded">
                                {{ currentBook.source }}
                            </div>

                            <!-- 竖向小图轮播（仅当有多本书时显示） -->
                            <div v-if="mainBooks.length > 1" class="absolute top-0 right-0 flex flex-col">
                                <button v-for="(book, index) in mainBooks" :key="`thumb-${index}`"
                                    @mouseenter="setCurrentIndex(index)"
                                    class="w-12 h-16 border-2 rounded-none overflow-hidden transition-all"
                                    :class="currentIndex === index ? 'border-primary' : 'border-white opacity-70 hover:opacity-100'">
                                    <img :src="book.cover" :alt="book.title" class="w-full h-full object-cover" />
                                </button>
                            </div>
                        </div>

                        <!-- 右侧：书籍信息 -->
                        <div class="w-full md:w-2/3 p-6">
                            <h2 class="text-[clamp(1.2rem,3vw,1.8rem)] font-bold mb-2 text-shadow cursor-pointer">
                                {{ currentBook.title }}
                            </h2>
                            <p class="text-gray-600 mb-4 line-clamp-7">
                                {{ currentBook.description }}
                            </p>
                            <div class="mb-4">
                                <span class="text-sm text-gray-500">简介：</span>
                                <span class="text-sm">{{ currentBook.tagline || '暂无简介' }}</span>
                            </div>
                            <a href="#"
                                class="inline-block bg-primary hover:bg-dark text-white px-4 py-2 rounded-md transition-colors cursor-pointer" @click.stop="navigateToBook(currentBook.id)">
                                开始阅读
                            </a>
                        </div>
                    </div>
                </div>

                <!-- 本周强推 -->
                <div class="w-full md:w-1/3">
                    <RankList title="本周强推" :books="weeklyBooks" />
                </div>
            </div>
        </section>

        <!-- 热门推荐和点击榜单 -->
        <section class="mb-10 w-5xl mx-auto">
            <div class="flex flex-col md:flex-row gap-6">
                <!-- 热门推荐 -->
                <div class="w-full md:w-2/3">
                    <h3 class="text-lg font-bold mb-4 border-b-2 border-primary pb-1 inline-block">热门推荐</h3>
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div v-for="(book, index) in recommendBooks" :key="index"
                            class="bg-white rounded-xl shadow-sm p-4 card-hover cursor-pointer" @click="navigateToBook(book.id)">
                            <div class="flex gap-4">
                                <img :src="book.cover" :alt="book.title" class="w-20 h-28 object-cover rounded" />
                                <div>
                                    <h4 class="font-medium">{{ book.title }}</h4>
                                    <p class="text-sm text-gray-600 line-clamp-2 mt-1">
                                        {{ book.description }}
                                    </p>
                                    <p class="text-xs text-gray-500 mt-2">作者：{{ book.author }}</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <!-- 点击榜单 -->
                <div class="w-full md:w-1/3">
                    <RankList title="点击榜单" :books="clickBooks" :show-more="true" more-link="/rank#clickest" />
                </div>
            </div>
        </section>

        <!-- 精品推荐和新书榜单 -->
        <section class="mb-10 w-5xl mx-auto">
            <div class="flex flex-col md:flex-row gap-6">
                <!-- 精品推荐 -->
                <div class="w-full md:w-2/3">
                    <h3 class="text-lg font-bold mb-4 border-b-2 border-primary pb-1 inline-block">精品推荐</h3>
                    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <div v-for="(book, index) in featuredBooks" :key="`featured-${index}`"
                            class="bg-white rounded-xl shadow-sm overflow-hidden card-hover cursor-pointer" @click="navigateToBook(book.id)">
                            <img :src="book.cover" :alt="book.title" class="w-full h-48 object-cover" />
                            <div class="p-3">
                                <h4 class="font-medium line-clamp-2">{{ book.title }}</h4>
                                <p class="text-xs text-gray-500 mt-1 line-clamp-2">
                                    {{ book.description }}
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
                <!-- 新书榜单 -->
                <div class="w-full md:w-1/3">
                    <RankList title="新书榜单" :books="newBooks" :show-more="true" more-link="/rank#newest" />
                </div>
            </div>
        </section>


        <!-- 最新更新和亚榜榜单 -->
        <section class="mb-10 w-5xl mx-auto">
            <div class="flex flex-col md:flex-row gap-6">
            <!-- 最新更新 -->
            <div class="w-full md:w-2/3">
                <h3 class="text-lg font-bold mb-4 border-b-2 border-primary pb-1 inline-block">最新更新</h3>
                <div class="bg-white rounded-xl shadow-sm overflow-hidden">
                <div class="overflow-x-auto">
                    <table class="w-full min-w-max">
                    <thead>
                        <tr class="bg-gray-50">
                        <th class="px-4 py-3 text-left text-sm font-semibold">类别</th>
                        <th class="px-4 py-3 text-left text-sm font-semibold">书名</th>
                        <th class="px-4 py-3 text-left text-sm font-semibold">最新章节</th>
                        <th class="px-4 py-3 text-left text-sm font-semibold">作者</th>
                        <th class="px-4 py-3 text-left text-sm font-semibold">更新时间</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr
                        v-for="(item, index) in latestBooks"
                        :key="index"
                        class=" hover:bg-gray-50 transition-colors cursor-pointer"
                        @click="navigateToBook(item.id)"
                        >
                        <td class="px-4 py-3 text-sm text-gray-600">{{ item.category }}</td>
                        <td class="px-4 py-3 text-sm font-medium">{{ item.title }}</td>
                        <td class="px-4 py-3 text-sm text-gray-600">{{ item.latestChapter }}</td>
                        <td class="px-4 py-3 text-sm text-gray-600">{{ item.author }}</td>
                        <td class="px-4 py-3 text-sm text-gray-500">{{ item.updateTime }}</td>
                        </tr>
                    </tbody>
                    </table>
                </div>
                </div>
            </div>

            <!-- 更新榜单  -->
            <div class="w-full md:w-1/3">
                <RankList title="更新榜单" :books="updateBooks" :show-more="true" more-link="/rank#updatest"/>
            </div>
            </div>
        </section>
    </div>

</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import Header from '@/components/Header.vue'
import RankList from '@/components/RankList.vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 当前显示第几本
const currentIndex = ref(0)

// 精品推荐轮播当前索引
const featuredCurrentIndex = ref(0)

// 导航到书籍详情页
const navigateToBook = (bookId) => {
    router.push(`/book/${bookId}`)
}


// 模拟多本书数据（后续可替换为 API）
const mainBooks = [
  {
    id: '1',
    title: '神医归来：十个女囚要我负责',
    description:
      '唐浩不想做神医，只想赚钱。从监狱释放后，六年了，竟有十个女囚找... 新产品：女人三枚就，唐浩不想做神医，只想赚钱。从监狱释放后，六年了，竟有十个女囚找...',
    tagline: '这个训练家太野了',
    cover: 'https://picsum.photos/id/24/500/700',
    source: '七猫中文网'
  },
  {
    id: '2',
    title: '赘婿：开局被退婚后我觉醒了',
    description:
      '林风穿越成赘婿，遭岳家羞辱退婚。觉醒系统后，他决定不再隐忍，一路逆袭...',
    tagline: '商战+爽文+打脸',
    cover: 'https://picsum.photos/id/25/500/700',
    source: '番茄小说'
  },
  {
    id: '3',
    title: '我在末世开超市',
    description:
      '丧尸爆发，别人在逃命，我却靠系统开了家超市。物资换忠诚，美女换积分...',
    tagline: '末世+种田+后宫',
    cover: 'https://picsum.photos/id/26/500/700',
    source: '飞卢小说网'
  }
]

// 计算属性：当前书籍
const currentBook = computed(() => mainBooks[currentIndex.value])

const weeklyBooks = [
  {
    id: '4',
    title: '都市：我不装了，反手拍碎少皇之墓',
    category: '抗日，谍战',
    quote: '"爹，我回来了，只是..."',
    cover: 'https://picsum.photos/id/20/100/140'
  },
  {
    id: '5',
    title: '穿书后，我在八十年代发家致富',
    category: '西方玄幻',
    quote: '万物从认知开始',
    cover: 'https://picsum.photos/id/21/100/140'
  },
  {
    id: '6',
    title: '开局：从万族战场证万人迷',
    category: '东方玄幻',
    quote: '我一个三绝剑帝十全能综合冠军吧',
    cover: 'https://picsum.photos/id/22/100/140'
  },
]

// 模拟数据（后续替换为 API 调用）
const recommendBooks = [
    {
        id: '10',
        title: '七零：靠弹棉花逆袭大佬爹多...',
        description: '穿成炮灰一天了，长兄是家中... 所有，我有一个大胆的想法，就是不做炮灰了，却被全家反对...',
        author: '小土豆',
        cover: 'https://picsum.photos/id/30/120/160'
    },
    {
        id: '11',
        title: '我娘子天下第一',
        description: '作 者：小一郎 标签：古风 武侠 搞笑 "娘，江南首富家公子哥来求亲..."',
        author: '小一郎',
        cover: 'https://picsum.photos/id/31/120/160'
    },
    {
        id: '12',
        title: '赘婿：我的前妻',
        description: '生活打的我措手不及，前妻跟别人跑了，我却一夜逆袭开始了我的复仇之路...',
        author: '爱了个寂寞QAQ',
        cover: 'https://picsum.photos/id/32/120/160'
    },
    {
        id: '13',
        title: '吾父即耶和华',
        description: '开局我妈告诉我个惊天秘密，我爸竟然是耶和华，耶稣是我哥，撒旦是我叔...',
        author: '天降猛男',
        cover: 'https://picsum.photos/id/33/120/160'
    },
    {
        id: '14',
        title: '重生了，回到小城当城门官',
        description: '大学毕业后，穷困潦倒了整五年，重活了回到了大学毕业，我只想搞钱买房...',
        author: '糊涂啊',
        cover: 'https://picsum.photos/id/34/120/160'
    },
    {
        id: '15',
        title: '快穿年代文妃',
        description: '【穿书+女强+甜宠+剧情+下乡】In 知识青年下乡，穿越年代的日子，顶着...',
        author: '拾光的灰烬',
        cover: 'https://picsum.photos/id/35/120/160'
    }
]

const clickBooks = [
    {
        id: '16',
        title: '不是主角，与隋唐罗成有什么关系',
        description: '李世民，历史原来不是这么念...',
        cover: 'https://picsum.photos/id/40/50/70'
    },
    {
        id: '17',
        title: '星河之上',
        description: '我在修仙界历经万古长青',
        cover: 'https://picsum.photos/id/41/50/70'
    },
    {
        id: '18',
        title: '六年后，我携四个幼崽翻墙回夫家',
        description: '我的七个美女姐姐的牛逼',
        cover: 'https://picsum.photos/id/42/50/70'
    },
    {
        id: '19',
        title: '龙王殿主，美女总裁求保证'
        // 无封面、无描述也支持
    },
    {
        id: '20',
        title: '独占仙胎'
    },
    {
        id: '21',
        title: '仙子别怕，我是厨子'
    },
    {
        id: '22',
        title: '大佬出狱，师妹我媳妇就在床'
    }
]

// 模拟数据（实际项目中可替换为 API 请求）
const featuredBooks = [
    {
        id: '23',
        title: '叛出家族后，转身投靠魔族女帝',
        description: '【玄幻+无系统+权谋+杀伐果断+情有独钟】宗门大比，他遭人暗算，修为尽废...',
        cover: 'https://picsum.photos/id/50/300/400'
    },
    {
        id: '24',
        title: '仙子别怕，我是厨子',
        description: '作者：赵公子 好好的仙门弟子，还得学做菜？没办法，谁让我只有一个大师傅呢...',
        cover: 'https://picsum.photos/id/51/300/400'
    },
    {
        id: '25',
        title: '六年后，我携四个幼崽翻墙',
        description: '作者：糖糖 因为误会，她离家不回，爹凉后，她带着四个缩小版的他霸气归来，惊艳众人...',
        cover: 'https://picsum.photos/id/52/300/400'
    },
    {
        id: '26',
        title: '傲太监：我乃大明九千岁',
        description: '【历史+权谋+腹黑+爽文】【伪太监】魂穿明朝，成为一名小太监，凭借聪明才智...',
        cover: 'https://picsum.photos/id/53/300/400'
    },
    {
        id: '27',
        title: '重生八零：我的毛概能刷出会模',
        description: '作者：橙子 上了大学，却发现毛概能刷出技能和物品...',
        cover: 'https://picsum.photos/id/54/300/400'
    },
    {
        id: '28',
        title: '乱世：从照顾嫂嫂开始修行',
        description: '一代剑神遭人陷害，重生归来，成了一名懵懂少年...',
        cover: 'https://picsum.photos/id/55/300/400'
    }
]

const newBooks = [
    {
        id: '29',
        title: '我愿为卿之实力者',
        description: '【玄幻爽文，男主天下第一】...',
        cover: 'https://picsum.photos/id/60/50/70'
    },
    {
        id: '30',
        title: '决裂：开局揭穿绿帽侠',
        description: '重生渔村，开局揭露老婆给我...',
        cover: 'https://picsum.photos/id/61/50/70'
    },
    {
        id: '31',
        title: '重生村医：开局回魂挽救全家',
        description: '大型狗血重生剧，瞬间泪目了',
        cover: 'https://picsum.photos/id/62/50/70'
    },
    { id: '32', title: '战神殿', description: '', cover: '' },
    { id: '33', title: '混世枭龙：从逐浪风帆开始', description: '', cover: '' },
    { id: '34', title: '七零：大佬，婶婶比你大很多', description: '', cover: '' }
]

// 模拟“最新更新”数据
const latestBooks = [
  {
    id: '35',
    category: '女生都市',
    title: '[HP]哈利与阿兹卡班',
    latestChapter: '371 Chapter 82',
    author: 'CisForCider',
    updateTime: '11/13 09:06'
  },
  {
    id: '36',
    category: '武侠修仙',
    title: '临时工物',
    latestChapter: '第431章：刀劈恶龙',
    author: '唐吉诃德',
    updateTime: '11/13 08:59'
  },
  {
    id: '37',
    category: '科幻灵异',
    title: '重生后，我被霸总反向追',
    latestChapter: '第139章 诡异的热搜',
    author: '我是江南',
    updateTime: '11/13 08:55'
  },
  {
    id: '38',
    category: '玄幻奇幻',
    title: '大奉打更人',
    latestChapter: '第1200章 杀许七安',
    author: '会说话的肘子',
    updateTime: '11/13 08:52'
  },
  {
    id: '39',
    category: '历史军事',
    title: '乱入三国之争霸召唤',
    latestChapter: '第187章惊天妙计定天化',
    author: '少年曹操',
    updateTime: '11/13 08:52'
  },
  {
    id: '40',
    category: '科幻灵异',
    title: '天医圣手',
    latestChapter: '第1660章 人生的路怎么走',
    author: '蓝色蝌蚪',
    updateTime: '11/13 08:08'
  }
]

// 模拟“更新榜单”数据
const updateBooks = [
  {
    id: '41',
    title: '[HP]哈利与阿兹卡班来亚',
    description: '在霍格沃茨的世界里留下的一个让人难忘的故事...',
    cover: 'https://picsum.photos/id/70/50/70'
  },
  {
    id: '42',
    title: '临时工物'
  },
  {
    id: '43',
    title: '重生后，我被霸总反向追'
  },
  {
    id: '44',
    title: '战神殿'
  },
  {
    id: '45',
    title: '乱入三国之争霸召唤'
  }
]

// 切换方法
function setCurrentIndex(index) {
  if (index >= 0 && index < mainBooks.length) {
    currentIndex.value = index
  }
}

onMounted(() => {
  // 精品推荐自动轮播
  const featuredTimer = setInterval(() => {
    featuredCurrentIndex.value = (featuredCurrentIndex.value + 1) % featuredBooks.length
  }, 3000)
    console.log('Component mounted.')
})
</script>

<style scoped>
.card-hover:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
  transition: all 0.3s ease;
  cursor: pointer;
}

.text-shadow {
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}
</style>