using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using famous.oauth;

namespace oauthdemo
{

  class Program
  {
    static string Base64Encode(string s)
    {
      var b = Encoding.UTF8.GetBytes(s);
      return Convert.ToBase64String(b);
    }
    private static readonly ClientSecrets BaiduSecrets = new ClientSecrets() { ClientId = "1X4E9T2cNAbGirXpOntbdYfN", ClientSecret = "0vmGPL773yYmH4zchYGg1ydpnNchoIM1" };
    private static readonly ClientSecrets WeiboSecrets = new ClientSecrets() { ClientId = "812692320", ClientSecret = "9b125acc087a3e372a7028be5bac053a" };
    // sohu should encode secret again
    private static readonly ClientSecrets SohuSecrets = new ClientSecrets() { ClientId = "36mk2EtAstxBqWpC9N6a", ClientSecret = Base64Encode("ZrMbr%JX^rxhrzhiyKt1o3j#JPhXs4d0Hyrjx=zs") };
    private static readonly ClientSecrets AzureSecrets = new ClientSecrets() { ClientId = "famous", ClientSecret = "SsdLfL8CH58qqGndYftP5w0jY75l+jFVkHxSJib9GpQ=" };
    private static readonly ClientSecrets NeteaseSecrets = new ClientSecrets() { ClientId = "LGBCrVyIcMe8prAE", ClientSecret = "POngu4hmQfRPxSuK2ALDNdv398e96sSZ" };
    
    static void Main(string[] args)
    {

      var ctx = new OAuth2Context(AzureConsts.AuthorizationServerUrl, AzureConsts.TokenServerUrl,
        AzureConsts.RedirectUrl)
      {
        ClientSecrets = AzureSecrets,
        DataStore = new FileDataStore("oauthdemo"),
//        Scopes = new[] { "https://api.datamarket.azure.com/" },// sohu require basic scope
      };
      var broker = new WebAuthorizationBroker(ctx);
      var resp = broker.AuthorizeAsync("default", true, CancellationToken.None).Result;
    }
  }
  static class BaiduConsts
  {
    public const string Base = "https://openapi.baidu.com/oauth/2.0/";
    public const string AuthorizationServerUrl = Base + "authorize";
    public const string TokenServerUrl = Base + "token";
    public const string RedirectUrl = "http://iweizhi2.duapp.com/authorize";
  }

  static class WeiboConsts
  {
    public const string Base = "https://api.weibo.com/oauth2/";
    public const string AuthorizationServerUrl = Base + "authorize";
    public const string TokenServerUrl = Base + "access_token";
    public const string RedirectUrl = "http://iweizhi2.duapp.com/authorize";
  }

  static class SohuConsts
  {
    public const string Base = "https://api.t.sohu.com/oauth2/";
    public const string AuthorizationServerUrl = Base + "authorize";
    public const string TokenServerUrl = Base + "access_token"; // https://api.t.sohu.com/oauth2/access_token
    public const string RedirectUrl = "http://iweizhi2.duapp.com/authorize";
    public const string Scopes = "basic";
  }
  static class NeteaseConsts
  {
    public const string Base = "https://api.t.163.com/oauth2/";
    public const string AuthorizationServerUrl = Base + "authorize";
    public const string TokenServerUrl = Base + "access_token"; // https://api.t.163.com/oauth2/access_token
    public const string RedirectUrl = "http://iweizhi2.duapp.com/authorize";
  }

  static class AzureConsts
  {
    public const string Base = "https://datamarket.azure.com/embedded/consent";
    public const string AuthorizationServerUrl = Base;
    public const string TokenServerUrl ="https://datamarket.accesscontrol.windows.net/v2/OAuth2-13"; 
    public const string RedirectUrl = "http://iweizhi2.duapp.com/authorize";
    public const string Scopes = "https://api.datamarket.azure.com/";
  }
}
