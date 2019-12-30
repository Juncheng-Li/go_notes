create 对象的时候可以使用Times关键字

```go
mockRepo.EXPECT().create(Any(), Any()).Return(nil).Times(5)


mockRepo.EXPECT().Retrieve(Any()).Return(objBytres1, nil)
mockRepo.EXPECT().Retrieve(Any()).Return(objBytres2, nil)
mockRepo.EXPECT().Retrieve(Any()).Return(objBytres3, nil)
mockRepo.EXPECT().Retrieve(Any()).Return(objBytres4, nil)
mockRepo.EXPECT().Retrieve(Any()).Return(objBytres5, nil)


retrieveCall := mockRepo.EXPECT().Retrieve(Any()).Return(nil, ErrAny)
mockRepo.EXPECT().create(Any(), Any()).Return(nil).After(retrieveCall)


InOrder(
    mockRepo.EXPECT().Retrieve(Any()).Return(nil, ErrAny)
    mockRepo.EXPECT().Create(Any(), Any()).Return(nil)
    mockRepo.EXPECT().Retrieve(Any()).Return(objBytes, nil)
)

func InOrder(calls ...*Call) {
    for i := 1; i < len(calls); i++ {
        calls[i].After(calls[i - 1])
    }
}
```