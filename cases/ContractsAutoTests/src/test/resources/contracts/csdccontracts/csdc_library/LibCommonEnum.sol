pragma solidity ^0.4.12;
/**
* @file LibCommonEnum.sol
* @author yiyating
* @time 2017-08-03
* @desc g等
*/


library LibCommonEnum {

    enum IdType{
        NONE,
        BISINESS_LICENSE,                           //1-营业执照
        CERTIFICATE_OF_REGISTRATION,                //2-登记证书
        ORGANIZATION_CODE_CERTIFICATE,              //3-组织机构代码证
        OFFICIAL_DOCUMENTS,                         //4-批文
        IDENTITY_CARD,                              //5-居民身份证
        PASSPORT,                                   //6-护照
        HOME_RETURN_PERMIT,                         //7-港澳居民来往内地通行证
        TAIWAN_RESIDENTS_PASS,                      //8-台湾居民往来大陆通行证
        FOREIGNERS_PERMANENT_PERMIT,                //9-外国人永久居留证
        HONGKONG_IDCARD,                            //10-香港居民身份证
        MACAO_IDCARD,                               //11-澳门居民身份证
        RESIDENCE_BOOKLET,                          //12-户口本
        SOCIAL_SECURITY,                            //13-社会保障号
        MILITARY_CARD,                              //14-军人证
        CIVILIAN_CARD,                              //15-文职证
        OFFICER_CARD,                               //16-警官证
        OTHER,                                      //17-其他证件
        UNIFIED_SOCIAL_CREDIT_CODE                  //18-统一社会信用代码
    }
}