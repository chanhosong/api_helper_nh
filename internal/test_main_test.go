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

package internal

import (
	"github.com/ghts/lib"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	f테스트_준비()
	테스트_실행결과 := m.Run()
	f테스트_정리()

	os.Exit(테스트_실행결과)
}

func f테스트_준비() {
	lib.F테스트_모드_시작()
	F_NH_API커넥터_실행()

	for i:=0 ; i<10 ; i++ {
		if  F접속됨_NH() {
			break
		}
	}

	lib.F조건부_패닉(!F접속됨_NH(), "접속 실패")

	var 에러 error
	소켓SUB_NH실시간_정보, 에러 = lib.New소켓SUB(lib.P주소_NH_실시간_CBOR)
	lib.F에러2패닉(에러)
}

func f테스트_정리() {
	lib.New소켓_질의(lib.P주소_NH_TR, lib.CBOR, lib.P10초).S질의(lib.S질의값_단순TR{TR구분:lib.TR종료})
	lib.F공통_종료_채널_닫은_후_재설정()
	lib.F에러2패닉(소켓SUB_NH실시간_정보.Close())
	lib.F테스트_모드_종료()
}
