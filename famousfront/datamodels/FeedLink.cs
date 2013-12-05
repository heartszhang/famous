using System.Runtime.Serialization;

namespace famousfront.datamodels
{
  [DataContract]
  internal class FeedLink
  {
    [DataMember(EmitDefaultValue = false)]
    public string uri
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string alias
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string local
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string cleaned_local
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public int words
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public int sentences
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public int links
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public int density
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public long length
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public bool readable
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public FeedMedia[] images
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public FeedMedia[] videos
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public FeedMedia[] audios
    {
      get;set;
    }
  }
}