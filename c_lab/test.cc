#include<iostream>
#include<gtest/gtest.h>
using namespace std;

int add(int a, int b) {
    return a+b;
}

TEST(test1, add) {
    EXPECT_EQ(add(1,3), 4);
}
int main(int argc, char *argv[]) {
    ::testing::InitGoogleTest(&argc, argv);
    return RUN_ALL_TESTS();
}
