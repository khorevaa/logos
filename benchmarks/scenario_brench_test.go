package benchmarks

import (
	"github.com/phuslu/log"
	"io/ioutil"

	//"io/ioutil"
	//"log"
	"testing"

	"go.uber.org/zap"
)

func BenchmarkDisabledWithoutFields(b *testing.B) {
	b.Logf("Logging at a disabled level without any structured context.")
	b.Run("Logos", func(b *testing.B) {
		logger := newZapLogger(zap.ErrorLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {

				logger.Info(getMessage(0))
			}
		})
	})
	//b.Run("Zap.Check", func(b *testing.B) {
	//	logger := newSampledLogger(zap.ErrorLevel)
	//	b.ResetTimer()
	//	b.RunParallel(func(pb *testing.PB) {
	//		for pb.Next() {
	//			if m := logger.Check(zap.InfoLevel, getMessage(0)); m != nil {
	//				m.Write()
	//			}
	//		}
	//	})
	//})
	//b.Run("Zap.Sugar", func(b *testing.B) {
	//	logger := newSampledLogger(zap.ErrorLevel).Sugar()
	//	b.ResetTimer()
	//	b.RunParallel(func(pb *testing.PB) {
	//		for pb.Next() {
	//			logger.Info(getMessage(0))
	//		}
	//	})
	//})
	//b.Run("Zap.SugarFormatting", func(b *testing.B) {
	//	logger := newZapLogger(zap.ErrorLevel).Sugar()
	//	b.ResetTimer()
	//	b.RunParallel(func(pb *testing.PB) {
	//		for pb.Next() {
	//			logger.Infof("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
	//		}
	//	})
	//})
	b.Run("apex/log", func(b *testing.B) {
		logger := newDisabledApexLog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newDisabledLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})

	b.Run("phuslu/log", func(b *testing.B) {

		logger := log.Logger{
			Level:      log.PanicLevel,
			TimeFormat: "", // uses rfc3339 by default
			Writer:     log.IOWriter{ioutil.Discard},
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg(getMessage(0))
			}
		})
	})
}

func BenchmarkAddingFields(b *testing.B) {
	b.Logf("Logging with additional context at each log site.")
	b.Run("logos", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Error(getMessage(0), fakeFields()...)
			}
		})
	})
	//b.Run("Zap.Check", func(b *testing.B) {
	//	logger := newSampledLogger(zap.DebugLevel)
	//	b.ResetTimer()
	//	b.RunParallel(func(pb *testing.PB) {
	//		for pb.Next() {
	//			if ce := logger.Check(zap.InfoLevel, getMessage(0)); ce != nil {
	//				//ce.Write(fakeFields()...)
	//			}
	//		}
	//	})
	//})
	//b.Run("Zap.CheckSampled", func(b *testing.B) {
	//	logger := newSampledLogger(zap.DebugLevel)
	//	b.ResetTimer()
	//	b.RunParallel(func(pb *testing.PB) {
	//		i := 0
	//		for pb.Next() {
	//			i++
	//			if ce := logger.Check(zap.InfoLevel, getMessage(i)); ce != nil {
	//				//ce.Write(fakeFields()...)
	//			}
	//		}
	//	})
	//})
	//b.Run("Zap.Sugar", func(b *testing.B) {
	//	logger := newSampledLogger(zap.DebugLevel).Sugar()
	//	b.ResetTimer()
	//	b.RunParallel(func(pb *testing.PB) {
	//		for pb.Next() {
	//			logger.Infow(getMessage(0), fakeSugarFields()...)
	//		}
	//	})
	//})
	b.Run("apex/log", func(b *testing.B) {
		logger := newApexLog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.WithFields(fakeApexFields()).Info(getMessage(0))
			}
		})
	})
	//b.Run("go-kit/kit/log", func(b *testing.B) {
	//	logger := newKitLog()
	//	b.ResetTimer()
	//	b.RunParallel(func(pb *testing.PB) {
	//		for pb.Next() {
	//			logger.Log(fakeSugarFields()...)
	//		}
	//	})
	//})
	//b.Run("inconshreveable/log15", func(b *testing.B) {
	//	logger := newLog15()
	//	b.ResetTimer()
	//	b.RunParallel(func(pb *testing.PB) {
	//		for pb.Next() {
	//			logger.Info(getMessage(0), fakeSugarFields()...)
	//		}
	//	})
	//})
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.WithFields(fakeLogrusFields()).Info(getMessage(0))
			}
		})
	})
	//b.Run("rs/zerolog", func(b *testing.B) {
	//	logger := newZerolog()
	//	b.ResetTimer()
	//	b.RunParallel(func(pb *testing.PB) {
	//		for pb.Next() {
	//			fakeZerologFields(logger.Info()).Msg(getMessage(0))
	//		}
	//	})
	//})
	b.Run("phuslu/log", func(b *testing.B) {

		logger := log.Logger{
			Level:      log.PanicLevel,
			TimeFormat: "", // uses rfc3339 by default
			Writer:     log.IOWriter{ioutil.Discard},
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {

				//logrus.Fields{
				//	"int":     _tenInts[0],
				//	"ints":    _tenInts,
				//	"string":  _tenStrings[0],
				//	"strings": _tenStrings,
				//	"time":    _tenTimes[0],
				//	"times":   _tenTimes,
				//	"user1":   _oneUser,
				//	"user2":   _oneUser,
				//	"users":   _tenUsers,
				//	"error":   errExample,
				//}
				logger.Info().
					Int("int", _tenInts[0]).
					Interface("ints", _tenInts).
					Str("string", _tenStrings[0]).
					Interface("strings", _tenStrings).
					Time("time", _tenTimes[0]).
					Interface("times", _tenTimes).
					Err(errExample).
					Msg(getMessage(0))

			}
		})
	})
}
