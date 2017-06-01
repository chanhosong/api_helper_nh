/* Copyright (C) 2015-2016 김운하(UnHa Kim)  unha.kim@kuh.pe.kr

이 파일은 GHTS의 일부입니다.

이 프로그램은 자유 소프트웨어입니다.
소프트웨어의 피양도자는 자유 소프트웨어 재단이 공표한 GNU LGPL 2.1판
규정에 따라 프로그램을 개작하거나 재배포할 수 있습니다.

이 프로그램은 유용하게 사용될 수 있으리라는 희망에서 배포되고 있지만,
특정한 목적에 적합하다거나, 이익을 안겨줄 수 있다는 묵시적인 보증을 포함한
어떠한 형태의 보증도 제공하지 않습니다.
보다 자세한 사항에 대해서는 GNU LGPL 2.1판을 참고하시기 바랍니다.
GNU LGPL 2.1판은 이 프로그램과 함께 제공됩니다.
만약, 이 문서가 누락되어 있다면 자유 소프트웨어 재단으로 문의하시기 바랍니다.
(자유 소프트웨어 재단 : Free Software Foundation, Inc.,
59 Temple Place - Suite 330, Boston, MA 02111-1307, USA)

Copyright (C) 2015년 UnHa Kim (unha.kim@kuh.pe.kr)

This file is part of GHTS.

GHTS is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, version 2.1 of the License.

GHTS is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with GHTS.  If not, see <http://www.gnu.org/licenses/>. */

package api_helper_nh

import (
	"github.com/ghts/lib"
	"github.com/go-mangos/mangos"
)

var 소켓SUB_NH실시간_정보 mangos.Socket
var 실시간_정보_중계_Go루틴_실행_중_MySQL = lib.New안전한_bool(false)
var nh구독채널_저장소_MySQL = lib.New구독채널_저장소()

func Go루틴_실시간_정보_중계(ch초기화 chan lib.T신호) {
	lib.F메모("모든 정보를 배포하는 대신, 채널별로 구독한 정보만 배포하는 방안을 생각해 볼 것.")

	if 실시간_정보_중계_Go루틴_실행_중.G값() {
		return
	}

	if 에러 := 실시간_정보_중계_Go루틴_실행_중.S값(true); 에러 != nil {
		ch초기화 <- lib.P신호_초기화
		return
	}

	defer 실시간_정보_중계_Go루틴_실행_중.S값(false)

	소켓SUB_NH실시간_정보, 에러 := lib.New소켓SUB(lib.P주소_NH_실시간_CBOR)
	lib.F에러2패닉(에러)

	nh대기_중_데이터_저장소 := new대기_중_데이터_저장소()
	ch종료 := lib.F공통_종료_채널()
	ch초기화 <- lib.P신호_초기화

	for {
		바이트_모음, 에러 := 소켓SUB_NH실시간_정보.Recv()
		if 에러 != nil {
			switch 에러.Error() {
			case "connection closed":
				소켓SUB_NH실시간_정보, 에러 = lib.New소켓SUB(lib.P주소_NH_실시간_CBOR)
			}

			lib.F에러_출력(에러)

			continue
		}

		소켓_메시지 := lib.New소켓_메시지from바이트_모음(바이트_모음)
		if 소켓_메시지.G에러() != nil {
			lib.F에러_출력(소켓_메시지.G에러())
			continue
		}

		for _, ch수신 := range nh구독채널_저장소.G채널_모음() {
			f실시간_정보_중계_도우미(ch수신, 소켓_메시지, nh구독채널_저장소, nh대기_중_데이터_저장소)
		}

		nh대기_중_데이터_저장소.S재전송()

		// 종료 조건 확인
		select {
		case <-ch종료:
			return
		default:    // 종료 신호 없으면 반복할 것.
		}

		lib.F실행권한_양보()
	}
}

// 이미 닫혀 있는 채널에 전송하면 패닉이 발생하는 데,
// 이 패닉을 적절하게 처리하도록 하기 위해서 해당 부분을 별도의 함수로 분리함.
func f실시간_정보_중계_도우미(ch수신 chan lib.I소켓_메시지, 소켓_메시지 lib.I소켓_메시지,
	nh구독채널_저장소 *lib.S구독채널_저장소, nh대기_중_데이터_저장소 *s대기_중_데이터_저장소) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{
		M함수with패닉내역: func(r interface{}) {
			lib.F에러_출력(r)
			// 이미 닫혀있는 채널에 전송하면 패닉 발생하므로, 해당 채널을 구독 채널 목록에서 제외시킴.
			nh구독채널_저장소.S삭제(ch수신)
		}})

	select {
	case ch수신 <- 소켓_메시지: // 메시지 중계 성공
	default: // 메시지 중계 실패.
		nh대기_중_데이터_저장소.S추가(소켓_메시지, ch수신)    // 저장소에 보관하고 추후 재전송
	}
}