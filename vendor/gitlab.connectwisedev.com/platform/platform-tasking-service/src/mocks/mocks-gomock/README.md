GoMock - a mock framework for Go.

Standard usage in our project:

0) (First time usage) To install binaries:

        go get github.com/golang/mock/gomock 
        cd $GOPATH/src/github.com/golang/mock/gomock 
        git checkout v1.1.1
        go install github.com/golang/mock/mockgen

1) Update the mocks using make command:
       
       cd src
       make generate-mock 
      
2) Use the mock in a test:

        func TestMyThing(t *testing.T) {
            testCase:= []struct{
                    mockTaskInstance mock.TaskInctanceConf
                }{
                    {
                        mockTaskInctance : func(ti *mock.MockTaskInstancePersistence) *mock.MockTaskInstancePersistence {
                            ti.EXPECT.SomeMethod("Only I Will Pass", gomock.Any).Return(errors.New(""))
                            return ti
                        }
                    }
                }
        
            for _, tc := range testCases {
            	tc := tc
            	t.Run(tc.name, func(t *testing.T) {
                    mockCtrl := gomock.NewController(t)
                    defer mockCtrl.Finish()
        
                    mockObj := something.NewMockMyInterface(mockCtrl)
                    if tc.mockTaskInctance !=nil{
                        mockObj = tc.mockTaskInstance(mockObj)
                    }
            
                    // pass mockObj to a real object using dependency injection
                    // or assign it to global variables and play with it.
                }
            }
        }
        
        
If test start panicing that means that some arguments that you pass to the function are wrong and you should change them        
        
Docks can be found [here](https://godoc.org/github.com/golang/mock/gomock)