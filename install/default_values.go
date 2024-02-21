package install

import (
	"github.com/dukhyungkim/gonuboard/model"
	"github.com/dukhyungkim/gonuboard/version"
	"path"
	"time"
)

var (
	defaultVersion        = version.Version
	defaultDataDirectory  = "data"
	defaultCacheDirectory = path.Join(defaultDataDirectory, "cache")
	defaultGrId           = "community"
	defaultReadPoint      = -1
	defaultWritePoint     = 5
	defaultCommentPoint   = 1
	defaultDownloadPoint  = -20

	defaultConfig = model.Config{
		CfID:                1,
		CfTitle:             defaultVersion,
		CfTheme:             "basic",
		CfAdminEmailName:    defaultVersion,
		CfUsePoint:          1,
		CfUseCopyLog:        1,
		CfLoginPoint:        100,
		CfCutName:           15,
		CfNickModify:        60,
		CfNewSkin:           "basic",
		CfNewRows:           15,
		CfSearchSkin:        "basic",
		CfConnectSkin:       "basic",
		CfFaqSkin:           "basic",
		CfReadPoint:         defaultReadPoint,
		CfWritePoint:        defaultWritePoint,
		CfCommentPoint:      defaultCommentPoint,
		CfDownloadPoint:     defaultDownloadPoint,
		CfWritePages:        10,
		CfMobilePages:       5,
		CfLinkTarget:        "_blank",
		CfDelaySec:          30,
		CfFilter:            "18아,18놈,18새끼,18뇬,18노,18것,18넘,개년,개놈,개뇬,개새,개색끼,개세끼,개세이,개쉐이,개쉑,개쉽,개시키,개자식,개좆,게색기,게색끼,광뇬,뇬,눈깔,뉘미럴,니귀미,니기미,니미,도촬,되질래,뒈져라,뒈진다,디져라,디진다,디질래,병쉰,병신,뻐큐,뻑큐,뽁큐,삐리넷,새꺄,쉬발,쉬밸,쉬팔,쉽알,스패킹,스팽,시벌,시부랄,시부럴,시부리,시불,시브랄,시팍,시팔,시펄,실밸,십8,십쌔,십창,싶알,쌉년,썅놈,쌔끼,쌩쑈,썅,써벌,썩을년,쎄꺄,쎄엑,쓰바,쓰발,쓰벌,쓰팔,씨8,씨댕,씨바,씨발,씨뱅,씨봉알,씨부랄,씨부럴,씨부렁,씨부리,씨불,씨브랄,씨빠,씨빨,씨뽀랄,씨팍,씨팔,씨펄,씹,아가리,아갈이,엄창,접년,잡놈,재랄,저주글,조까,조빠,조쟁이,조지냐,조진다,조질래,존나,존니,좀물,좁년,좃,좆,좇,쥐랄,쥐롤,쥬디,지랄,지럴,지롤,지미랄,쫍빱,凸,퍽큐,뻑큐,빠큐,ㅅㅂㄹㅁ",
		CfMemberSkin:        "basic",
		CfRegisterLevel:     2,
		CfRegisterPoint:     1000,
		CfIconLevel:         2,
		CfUseRecommend:      0,
		CfRecommendPoint:    0,
		CfLeaveDay:          30,
		CfSearchPart:        10000,
		CfEmailUse:          1,
		CfProhibitID:        "admin,administrator,관리자,운영자,어드민,주인장,webmaster,웹마스터,sysop,시삽,시샵,manager,매니저,메니저,root,루트,su,guest,방문객",
		CfProhibitEmail:     "",
		CfNewDel:            30,
		CfMemoDel:           180,
		CfVisitDel:          180,
		CfPopularDel:        180,
		CfUseMemberIcon:     2,
		CfMemberIconSize:    5000,
		CfMemberIconWidth:   22,
		CfMemberIconHeight:  22,
		CfMemberImgSize:     50000,
		CfMemberImgWidth:    60,
		CfMemberImgHeight:   60,
		CfLoginMinutes:      10,
		CfImageExtension:    "gif|jpg|jpeg|png|webp",
		CfFlashExtension:    "swf",
		CfMovieExtension:    "asx|asf|wmv|wma|mpg|mpeg|mov|avi|mp3",
		CfFormmailIsMember:  1,
		CfPageRows:          15,
		CfMobilePageRows:    15,
		CfStipulation:       "해당 홈페이지에 맞는 회원가입약관을 입력합니다.",
		CfPrivacy:           "해당 홈페이지에 맞는 개인정보처리방침을 입력합니다.",
		CfMobileNewSkin:     "basic",
		CfMobileSearchSkin:  "basic",
		CfMobileConnectSkin: "basic",
		CfMobileFaqSkin:     "basic",
		CfMobileMemberSkin:  "basic",
		CfCaptchaMp3:        "basic",
		CfEditor:            "ckeditor4",
		CfCertLimit:         2,
	}
	defaultMember = model.Member{
		MbLevel:        10,
		MbMailling:     1,
		MbOpen:         1,
		MbNickDate:     time.Now(),
		MbEmailCertify: time.Now(),
		MbDatetime:     time.Now(),
		MbIP:           "127.0.0.1",
	}
	defaultContents = []model.Content{
		{
			CoID:         "company",
			CoHTML:       1,
			CoSubject:    "회사소개",
			CoContent:    "<p align=center><b>회사소개에 대한 내용을 입력하십시오.</b></p>",
			CoSkin:       "basic",
			CoMobileSkin: "basic",
		},
		{
			CoID:         "provision",
			CoHTML:       1,
			CoSubject:    "서비스 이용약관",
			CoContent:    "<p align=center><b>서비스 이용약관에 대한 내용을 입력하십시오.</b></p>",
			CoSkin:       "basic",
			CoMobileSkin: "basic",
		},
		{
			CoID:         "privacy",
			CoHTML:       1,
			CoSubject:    "개인정보 처리방침",
			CoContent:    "<p align=center><b>개인정보 처리방침에 대한 내용을 입력하십시오.</b></p>",
			CoSkin:       "basic",
			CoMobileSkin: "basic",
		},
	}
	defaultQAConfig = model.QaConfig{
		QaTitle:            "1:1문의",
		QaCategory:         "회원|포인트",
		QaSkin:             "basic",
		QaMobileSkin:       "basic",
		QaUseEmail:         1,
		QaReqEmail:         0,
		QaUseHp:            1,
		QaReqHp:            0,
		QaUseEditor:        1,
		QaSubjectLen:       60,
		QaMobileSubjectLen: 30,
		QaPageRows:         15,
		QaMobilePageRows:   15,
		QaImageWidth:       600,
		QaUploadSize:       1048576,
		QaInsertContent:    "",
	}
	defaultGroup = model.Group{
		GrID:      defaultGrId,
		GrSubject: "커뮤니티",
	}
)
