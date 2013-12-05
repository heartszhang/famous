using System.Runtime.Serialization;

namespace famousfront.datamodels
{
  [DataContract]
  internal class FeedEntry
  {
    [DataMember(Name = "_id")]
    public string id
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public uint flags
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string source
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public uint type
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string uri
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public FeedTitle title
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public FeedAuthor author
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public long pubdate
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string summary
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string content
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string[]tags
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public FeedMedia[] images
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public FeedMedia[]videos
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public FeedMedia[]audios
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public FeedLink[] links
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public uint words
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public uint density
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public ulong status
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string[] categories
    {
      get;set;
    }
  }
}