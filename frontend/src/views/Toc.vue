<template>
    <!-- 小屏幕 -->
    <div :class="['block md:hidden min-h-screen', currentTheme === 'night' ? 'night' : '', currentTheme === 'eye-protect' ? 'eye-protect' : '']">
        <!-- 顶部导航栏 -->
        <header class="sticky top-0 z-50 bg-[#469b75] text-white shadow-md">
            <div class="flex items-center justify-between px-4 py-3">
                <button @click="goBack" class="text-white text-xl cursor-pointer">
                    <i class="fas fa-arrow-left"></i>
                </button>
                <h1 class="text-base font-bold truncate mx-2">{{ chapterTitle }}</h1>
                <button @click="goHome" class="text-white text-xl cursor-pointer">
                    <i class="fas fa-home"></i>
                </button>
            </div>
        </header>
        
        <!-- 阅读控制栏 -->
        <div :class="[
            'sticky top-[56px] z-40 flex items-center justify-between px-4 py-2 theme-border transition-transform duration-300',
            showControls ? 'translate-y-0' : '-translate-y-full',
            'theme-control-bg'
        ]">
            <div class="flex items-center gap-2">
                <button 
                    @click="toggleNightMode"
                    :class="['px-3 py-1 rounded-full text-sm font-medium shadow-sm transition-all', 
                            currentTheme === 'night' ? 'bg-gradient-to-r from-purple-600 to-purple-700 text-white' : 'bg-gradient-to-r from-gray-300 to-gray-400 text-gray-700 hover:from-gray-400 hover:to-gray-500']"
                >
                    关灯
                </button>
                <button 
                    @click="toggleEyeProtect"
                    :class="['px-3 py-1 rounded-full text-sm font-medium shadow-sm transition-all', 
                            currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-500 to-green-600 text-white' : 'bg-gradient-to-r from-gray-300 to-gray-400 text-gray-700 hover:from-gray-400 hover:to-gray-500']"
                >
                    护眼
                </button>
            </div>
            <div class="flex items-center gap-2">
                <span class="text-sm theme-text">字体：</span>
                <button 
                    class="px-2 py-1 rounded-full text-sm font-medium shadow-sm transition-all"
                    :class="currentTheme === 'night' ? 'bg-gradient-to-r from-gray-700 to-gray-800 text-white border border-gray-600' : 
                            currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-100 to-green-200 text-green-800 border border-green-300' : 
                            'bg-gradient-to-r from-white to-gray-100 text-gray-700 border border-gray-200 hover:from-gray-100 hover:to-gray-200'"
                >大</button>
                <button 
                    class="px-2 py-1 rounded-full text-sm font-medium shadow-sm transition-all"
                    :class="currentTheme === 'night' ? 'bg-gradient-to-r from-purple-600 to-purple-700 text-white border border-purple-500' : 
                            currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-500 to-green-600 text-white border border-green-400' : 
                            'bg-gradient-to-r from-[#469b75] to-[#3d8766] text-white border border-[#3d8766]'"
                >中</button>
                <button 
                    class="px-2 py-1 rounded-full text-sm font-medium shadow-sm transition-all"
                    :class="currentTheme === 'night' ? 'bg-gradient-to-r from-gray-700 to-gray-800 text-white border border-gray-600' : 
                            currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-100 to-green-200 text-green-800 border border-green-300' : 
                            'bg-gradient-to-r from-white to-gray-100 text-gray-700 border border-gray-200 hover:from-gray-100 hover:to-gray-200'"
                >小</button>
            </div>
            <button 
                class="w-10 h-10 flex items-center justify-center rounded-full shadow-sm transition-all"
                :class="currentTheme === 'night' ? 'bg-gradient-to-r from-gray-700 to-gray-800 text-white' : 
                        currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-400 to-green-500 text-white' : 
                        'bg-gradient-to-r from-gray-300 to-gray-400 text-gray-700 hover:from-gray-400 hover:to-gray-500'"
            >
                <i class="fas fa-volume-up text-lg"></i>
            </button>
        </div>
        
        <!-- 直达底部按钮 -->
        <div class="fixed right-4 bottom-20 z-30">
            <button @click="scrollToBottom" 
                class="w-12 h-12 bg-[#469b75] text-white rounded-full shadow-lg flex items-center justify-center hover:bg-[#3d8766] transition-colors"
                aria-label="直达底部">
                <i class="fas fa-arrow-down text-lg"></i>
            </button>
        </div>
        
        <!-- 章节列表 -->
        <div class="p-4 ">
            <div :class="['rounded-lg shadow-md p-4 mb-4', 
                        currentTheme === 'night' ? 'bg-[#2d2d2d] theme-border' : 
                        currentTheme === 'eye-protect' ? 'bg-[#b3e5cc] theme-border' : 'bg-white']">
                <h2 :class="['text-lg font-semibold mb-4', 
                          currentTheme === 'night' ? 'text-white' : 
                          currentTheme === 'eye-protect' ? 'text-[#333333]' : 'text-gray-800']">
                    正文({{ chapters.length }})
                </h2>
                <div class="space-y-2">
                    <a v-for="(chapter, index) in chapters" :key="index" @click="goToChapter(index + 1, chapter)"
                        :class="['flex items-center justify-between cursor-pointer px-3 py-3 rounded transition-colors border', 
                                  currentTheme === 'night' ? 
                                  'text-white hover:text-emerald-400 hover:bg-[#3d3d3d] border-[#404040]' : 
                                  currentTheme === 'eye-protect' ? 
                                  'text-[#333333] hover:text-[#3d8766] hover:bg-[#a5d6a7] border-[#81c784]' : 
                                  'text-gray-700 hover:text-emerald-600 hover:bg-gray-50 border-gray-100']">
                        <span class="flex-grow truncate pr-2">{{ chapter }}</span>
                        <span class="text-emerald-600 text-sm min-w-12 text-right">[免费]</span>
                    </a>
                </div>
            </div>
        </div>

        <!-- 底部导航栏 -->
        <ToolBar :showControls="showControls" />
    </div>

    <!-- 大屏 -->
    <div class="hidden md:block container md:w-5xl  mx-auto ">
        <div class="hidden md:block">
            <Header />
        </div>
        <!-- Breadcrumb -->
        <div class="bg-gray-100 ">
            <div class="max-w-7xl mx-auto px-4 py-3 text-sm text-gray-600">
                <a href="#" class="hover:text-emerald-600">小说精品屋</a>
                <span class="mx-2">></span>
                <a href="#" class="hover:text-emerald-600">{{ bookData.breadcrumb.category }}</a>
                <span class="mx-2">></span>
                <a href="#" class="hover:text-emerald-600">{{ bookData.title }}</a>
                <span class="mx-2">></span>
                <span>作品目录</span>
            </div>
        </div>

        <!-- Main Content -->
        <main class="max-w-7xl mx-auto px-4 py-8 bg-white ">
            <h1 class="text-3xl font-bold text-center text-gray-800 mb-6">
                {{ bookData.title }}
            </h1>

            <div class="flex justify-center gap-8 text-sm text-gray-600 mb-8">
                <span>作者：{{ bookData.author }}</span>
                <span>分类：{{ bookData.breadcrumb.category }}</span>
                <span>状态：{{ bookData.status }}</span>
                <span>总点击：{{ bookData.totalClicks }}</span>
                <span>总字数：{{ bookData.totalWords }}</span>
            </div>

            <h2 class="text-xl font-semibold text-gray-800 bg-gray-100 mb-6">正文(100)</h2>

            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-x-6 gap-y-3">
                <a v-for="(chapter, index) in chapters" :key="index" @click="goToChapter(index + 1, chapter)"
                    :class="['flex items-center justify-between cursor-pointer px-3 py-2 rounded transition-colors', 
                              currentTheme === 'night' ? 
                              'text-white hover:text-emerald-400 hover:bg-gray-800' : 
                              currentTheme === 'eye-protect' ? 
                              'text-[#333333] hover:text-[#3d8766] hover:bg-[#a5d6a7]' : 
                              'text-gray-700 hover:text-emerald-600 hover:bg-gray-50']">
                    <span class="flex-grow truncate pr-2">{{ chapter }}</span>
                    <span class="text-emerald-600 text-sm min-w-16 text-right">[免费]</span>
                </a>
            </div>
        </main>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted  } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Header from '@/components/Header.vue';
import ToolBar from '@/components/ToolBar.vue';

const route = useRoute();
const router = useRouter();

// 小屏幕相关状态
const currentTheme = ref('default');
const chapterTitle = ref('');
const showControls = ref(true);
const lastScrollTop = ref(0);
const scrollThreshold = 10;

// 书籍数据
const bookData = ref({
    title: '',
    author: '',
    category: '',
    status: '',
    totalClicks: 0,
    totalWords: 0,
    breadcrumb: {
        category: ''
    }
});

const chapters = ref([]);

// 小屏幕相关方法
const goBack = () => {
    router.back();
};

const goHome = () => {
    router.push('/');
};

const toggleNightMode = () => {
    currentTheme.value = currentTheme.value === 'night' ? 'default' : 'night';
};

const toggleEyeProtect = () => {
    currentTheme.value = currentTheme.value === 'eye-protect' ? 'default' : 'eye-protect';
};

const scrollToBottom = () => {
    window.scrollTo({
        top: document.body.scrollHeight,
        behavior: 'smooth'
    });
};

// 跳转到章节阅读页面
const goToChapter = (chapterIndex, chapterTitle) => {
    // 构建模拟的chapterId，实际项目中应该从API获取
    const chapterId = `${chapterIndex}`;
    router.push(`/book/${bookData.value.bookId}/${chapterId}`);
};

// 滚动监听处理函数
const handleScroll = () => {
    const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
    
    // 只有当滚动超过阈值时才触发显示/隐藏逻辑
    if (Math.abs(scrollTop - lastScrollTop.value) > scrollThreshold) {
        // 向上滚动时显示控制栏，向下滚动时隐藏
        showControls.value = scrollTop < lastScrollTop.value;
        lastScrollTop.value = scrollTop;
    }
};

// 获取书籍数据
onMounted(() => {
    // 添加滚动监听
    window.addEventListener('scroll', handleScroll);
    // 获取路由参数
    const bookId = route.params.bookId;

    // 模拟数据，实际项目中应该从API获取
    bookData.value = {
        bookId: bookId,
        title: '七零：靠弹幕嫁绝嗣大佬多胎啦！',
        author: '我莫波顿了',
        category: '女生频道',
        status: '连载中',
        totalClicks: 1456,
        totalWords: 146919,
        breadcrumb: {
            home: '小说精品屋',
            category: '女生频道'
        }
    };
    
    // 设置小屏幕导航栏标题
    chapterTitle.value = bookData.value.title;

    // 模拟目录数据
    chapters.value = [
        "第一卷 第1章 看见奇怪的文字",
        "第一卷 第2章 白月光宝贝",
        "第一卷 第3章 得知身世，修金手指",
        "第一卷 第4章 知晓郝家人的罪证",
        "第一卷 第5章 开介绍信，你为什么要帮我？",
        "第一卷 第6章 洛阳王，学习",
        "第一卷 第7章 滚男要办酒席，十个罪可的散漫",
        "第一卷 第8章 调查郑家，这口味挺好……另类",
        "第一卷 第9章 举门躲给她老汉",
        "第一卷 第10章 呸，败人，你败",
        "第一卷 第11章 调查，郑家你事便区",
        "第一卷 第12章 找要听家",
        "第一卷 第13章 还是……你在她我们？",
        "第一卷 第14章 下药，你的胎呢？",
        "第一卷 第15章 死了，卸虎完了",
        "第一卷 第16章 卫撑埠光",
        "第一卷 第17章 温我去光",
        "第一卷 第18章 举报你",
        "第一卷 第19章 退给朴",
        "第一卷 第20章 被调查，病倒",
        "第一卷 第21章 商爸攻的誓言，洼男传师将",
        "第一卷 第22章 中风，出发",
        "第一卷 第23章 抱位置，金容",
        "第一卷 第24章 如就，修攻你播",
        "第一卷 第25章 酒洒，大怕母大剛",
        "第一卷 第26章 叶母的门难，名声严",
        "第一卷 第27章 醉杯，抵达海岛",
        "第一卷 第28章 皇帝，又工团招到",
        "第一卷 第29章 考试，收怡宿舍",
        "第一卷 第30章 宿友，八挂",
        "第一卷 第31章 邀请，把影您扣人头上了",
        "第一卷 第32章 精周，加位显韩胖",
        "第一卷 第33章 幼般，秘密任务",
        "第一卷 第34章 联这合，私下找家件合作",
        "第一卷 第35章 钱票，谈心",
        "第一卷 第36章 打架，真的好棒！",
        "第一卷 第37章 打胸，小姑不行老虎上",
        "第一卷 第38章 失望，拿命道题",
        "第一卷 第39章 天威强，想挖人",
        "第一卷 第40章 下毒，食莹倒到",
        "第一卷 第41章 反击，提醒，只要居团长肯放人",
        "第一卷 第42章 告白，男二，如他斯来的恶意",
        "第一卷 第43章 她想要死我！",
        "第一卷 第44章 流言起，各种攻意到她茶来",
        "第一卷 第45章 找持，制你对方自便",
        "第一卷 第46章 为什么要这样对好不容易得来的",
        "第一卷 第47章 难道她就不配店吾吗？",
        "第一卷 第48章 食人的岱底",
        "第一卷 第49章 找您给天才",
        "第一卷 第50章 他是在为他死去的采妹姆罪",
        "第一卷 第51章 怎么也攻不了要死的结局吗？",
        "第一卷 第52章 得救，有这么多人关心她",
        "第一卷 第53章 下乡，你是想就此报告",
        "第一卷 第54章 成为律师，被秋为孙子",
        "第一卷 第55章 两夜肃开",
        "第一卷 第56章 拦路中，两攻行走",
        "第一卷 第57章 生病，有人被强吗",
        "第一卷 第58章 拿刀相扣！",
        "第一卷 第59章 瘦幻，把刀给我",
        "第一卷 第60章 你不要命了！",
        "第一卷 第61章 求你了，帮帮我",
        "第一卷 第62章 野菜树根，被楼上",
        "第一卷 第63章 西厨房",
        "第一卷 第64章 全村人道证",
        "第一卷 第65章 村长卖徒是特务头子？",
        "第一卷 第66章 李如青，求你了，放开吧！",
        "第一卷 第67章 为什么不能映点，再试一次",
        "第一卷 第68章 实死她了，始始吃",
        "第一卷 第69章 只要你愿能君，我愿你做",
        "第一卷 第70章 半夜大火",
        "第一卷 第71章 逃出",
        "第一卷 第72章 发财了！",
        "第一卷 第73章 卓家",
        "第一卷 第74章 坦象心声，找媳你",
        "第一卷 第75章 当后妈，疾了",
        "第一卷 第76章 我没有你个小缺就",
        "第一卷 第77章 就为给自己个公道",
        "第一卷 第78章 怎么，不欢迎吗",
        "第一卷 第79章 妈，你怎郝死我吗",
        "第一卷 第80章 正主都来了，怎么，还舍不得归",
        "第一卷 第81章 示羽，证不会",
        "第一卷 第82章 你，你怎扣我家人了？",
        "第一卷 第83章 得知小缺就是郑听开",
        "第一卷 第84章 拿，你给我留点面面",
        "第一卷 第85章 每个人的人生不同",
        "第一卷 第86章 我和郑同志愿说人关系",
        "第一卷 第87章 清宁原读你吗",
        "第一卷 第88章 初见，清草",
        "第一卷 第89章 不，是不好容",
        "第一卷 第90章 清草",
        "第一卷 第91章 发疯的郑斯斯",
        "第一卷 第92章 坦白",
        "第一卷 第93章 愚南",
        "第一卷 第94章 卓老爷子发威",
        "第一卷 第95章 固埋失去车山，修赦",
        "第一卷 第96章 不，那是要她的命",
        "第一卷 第97章 不然，免谈！",
        "第一卷 第98章 过世",
        "第一卷 第99章 你是我对象",
        "第一卷 第100章 大结局"
    ];
});

// 组件卸载时移除滚动监听
onUnmounted(() => {
    window.removeEventListener('scroll', handleScroll);
});
</script>

<style scoped>
/* 可选：添加局部样式 */
</style>