// Copyright 2018-2019 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package snapshotdb

/*
func TestCloseJournalWriter(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "test_close*.log")
	if err != nil {
		t.Error(err)
	}
	jw := newJournalWriter(f)
	writer, err := jw.journal.Next()
	if err != nil {
		t.Error(err)
	}
	if _, err := writer.Write([]byte("a")); err != nil {
		t.Error("should write", err)
	}
	if err := jw.Close(); err != nil {
		t.Error("should can close", err)
	}
	if _, err := jw.journal.Next(); err == nil {
		t.Fatal(err)
	}
	if err := jw.writer.Close(); err == nil {
		t.Error("should have be closed")
	}
}
*/
//
//func TestRMthan(t *testing.T) {
//	Instance()
//	if err := dbInstance.writeJournalHeader(big.NewInt(1), generateHash("aaaa"), common.ZeroHash, journalHeaderFromUnRecognized); err != nil {
//		t.Error(err)
//	}
//	if err := dbInstance.writeJournalBody(generateHash("aaaa"), []byte("abcdefg")); err != nil {
//		t.Error(err)
//	}
//	fds, err := dbInstance.storage.List(TypeJournal)
//	if err != nil {
//		t.Error(err)
//	}
//
//	for _, fd := range fds {
//		reader, err := dbInstance.storage.Open(fd)
//		if err != nil {
//			t.Error(err)
//		}
//		journals := journal.NewReader(reader, nil, false, false)
//		j, err := journals.Next()
//		if err != nil {
//			t.Error(err)
//		}
//		var header journalHeader
//		if err := decode(j, &header); err != nil {
//			t.Error(err)
//		}
//		logger.Debug("header", "v", header)
//		for {
//			j, err := journals.Next()
//			if err == io.EOF {
//				break
//			}
//			if err != nil {
//				t.Error(err)
//			}
//			var body string
//			if err := decode(j, &body); err != nil {
//				t.Error(err)
//			}
//			logger.Debug("body", "v", body, "byte", j)
//		}
//	}
//
//	if err := dbInstance.writeJournalHeader(big.NewInt(1), generateHash("aaaa"), common.ZeroHash, journalHeaderFromRecognized); err != nil {
//		t.Error(err)
//	}
//	if err := dbInstance.writeJournalBody(generateHash("aaaa"), []byte("bbbbbbbbbbbbbbbbbbbb")); err != nil {
//		t.Error(err)
//	}
//
//	for _, fd := range fds {
//		reader, err := dbInstance.storage.Open(fd)
//		if err != nil {
//			t.Error(err)
//		}
//		journals := journal.NewReader(reader, nil, false, false)
//		j, err := journals.Next()
//		if err != nil {
//			t.Error(err)
//		}
//		var header journalHeader
//		if err := decode(j, &header); err != nil {
//			t.Error(err)
//		}
//		logger.Debug("header", "v", header)
//		for {
//			j, err := journals.Next()
//			if err == io.EOF {
//				break
//			}
//			if err != nil {
//				t.Error(err)
//			}
//			var body string
//			if err := decode(j, &body); err != nil {
//				t.Error(err)
//			}
//			logger.Debug("body", "v", body, "byte", j)
//		}
//	}
//}
