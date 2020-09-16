#!/bin/bash

echo "begin test dns proxy!"

n=$1
dnsProxy=$2
dnsTrans=$3

searchDomain(){
    timer_start=`date "+%s"`

    iter=1
    while(( $iter<=$1 ))
    do
        nslookup baidu.com $2
        let "iter++"
    done

    timer_end=`date "+%s"`

    return $(($timer_end-$timer_start))
}

searchDomain $n $dnsProxy
timePass=$?

searchDomain $n $dnsTrans
timePassTrans=$?

searchDomain $n 
timePassDefault=$?

echo "查询DNS中继服务 查询( $n )次 花费时间( $timePass )秒"
echo "查询转发服务$dnsTrans 查询( $n )次 花费时间( $timePassTrans )秒"
echo "查询默认DNS服务 查询( $n )次 花费时间( $timePassDefault )秒"