package hw03frequencyanalysis

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// русский текст.
var ruText = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

var expectedRuStringSet = "а он и ты что в его если кристофер не"

// английский текст.
var enText = `Go provides two features that replace class inheritance.[citation needed]
The first is embedding, which can be viewed as an automated form of composition.[70]
The second are its interfaces, which provides runtime polymorphism.
Interfaces are a class of types and provide a limited form of structural typing in the 
otherwise nominal type system of Go. An object which is of an interface type is also of 
another type, much like C++ objects being simultaneously of a base and derived class. 
Go interfaces were designed after protocols from the Smalltalk programming language. 
Multiple sources use the term duck typing when describing Go interfaces. 
Although the term duck typing is not precisely defined and therefore not wrong, 
it usually implies that type conformance is not statically checked. Because conformance 
to a Go interface is checked statically by the Go compiler (except when performing a type assertion), 
the Go authors prefer the term structural typing.`

var expectedEnStringSet = "the go of is a type interfaces typing an and"

// текст на санскрите (не знаю что там. что то при цветы и фрукты. взято из сети для теста )).
var sanskritTest = `कुसुमानि च फलानि च
देवान् कुसुमैः पूजयन्ति नराः 
फलैर् अपि पूजयन्ति 
कुसुमानि च फलानि च वनाद् आनयामि 
देवदत्तेन सह तत्र गच्छामि 
कुसुमानि च फलानि च वनस्य वृक्षेषु रोहन्ति`
var expectedSanskritStringSet = "च कुसुमानि फलानि पूजयन्ति अपि आनयामि कुसुमैः गच्छामि तत्र देवदत्तेन"

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("string with punctuation symbols only", func(t *testing.T) {
		require.Len(t, Top10(".,!?'\"/()\\[]{}–"), 0)
	})

	t.Run("string with word count < 10", func(t *testing.T) {
		require.Len(t, Top10("test for fore fore words"), 4)
	})

	t.Run("positive russian test", func(t *testing.T) {
		expected := strings.Fields(expectedRuStringSet)
		require.Equal(t, expected, Top10(ruText))
	})

	t.Run("positive english test", func(t *testing.T) {
		expected := strings.Fields(expectedEnStringSet)
		require.Equal(t, expected, Top10(enText))
	})

	t.Run("positive sanskrit test", func(t *testing.T) {
		expected := strings.Fields(expectedSanskritStringSet)
		require.Equal(t, expected, Top10(sanskritTest))
	})
}
