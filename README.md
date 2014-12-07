
--> GTF MUST be put under $gopath\src, NOT under $gopath\src\github.com\ZhangQin3 <--

This is a generic test automation framework based on golang, it is named as Generic Testing Framework(GTF). Now, GTF is in its very early stage, but it has provide most necessary features that a mature testing framework shall provide, it can be used in unit test, integration test, system test and acceptance test with only a little simple adaption.

With GTF, one test script is a golang package in a single .go file, all the test scripts are put in the gtf/scripts directory. Related test scripts can get together and are put into a test suite, the test suite is also an .go file and a golang package located in the gtf/testsuites directory.

GTF generates one executable file(.exe) for each test suite, the executable file is put in the $gopath/bin/ directory. To generate an executable file for a test suite, you should build and run ignition.go package in gtf/drivers/ignition directory, this will generate a .exe file and puts it in $gopath/bin/.

Run the executable file( execute.exe for now) in $gopath/bin/, all the tests related with the test suite will be executed and logs for each test scripts will be generated and put in $gopath/bin/.

GTF Provides:
1) Test case level and test script level and test suite level setup and teardown.
2) HTML log output for each test script and the log file gives logs for each test case.
3) A simple way for creating customized test libraries(such as, Selenium, Telnet, SSH) with golang std libraries and gtf/log.


TBD:
data driven/keyword driven support. 
log page internal jump. 

Only support windows now.
